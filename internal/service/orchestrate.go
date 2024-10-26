package service

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gambtho/zkillanalytics/internal/api/zkill"
	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/utils"
)

// GetAllData Main function to orchestrate data fetching based on availability and necessity.
func GetAllData(client *http.Client, corporations, alliances, characters []int, startString, endString string) (*model.ChartData, error) {
	fmt.Println(fmt.Sprintf("Fetching data for %s to %s...", startString, endString))
	fetchStart := time.Now()

	startDate, _ := time.Parse("2006-01-02", startString)
	endDate, _ := time.Parse("2006-01-02", endString)
	year := startDate.Year()
	esiRefresh := false

	dataAvailability, _ := utils.CheckDataAvailability(int(startDate.Month()), int(endDate.Month()), year)

	esiFileName := persist.GenerateEsiDataFileName()
	fileInfo, err := os.Stat(esiFileName)

	esiData, err := persist.ReadEsiDataFromFile(esiFileName)
	if err != nil || time.Since(fileInfo.ModTime()) > 48*time.Hour || fileInfo.Size() <= 50*1024 {
		fmt.Println("Error loading esi data from file:", err)
		esiData = &model.ESIData{
			AllianceInfos:    make(map[int]model.Alliance),
			CharacterInfos:   make(map[int]model.Character),
			CorporationInfos: make(map[int]model.Corporation),
		}
		esiRefresh = true
	}

	fetchIDs := &model.Ids{
		CorporationIDs: corporations,
		AllianceIDs:    alliances,
		CharacterIDs:   characters,
	}

	idChanged, newIDs, err := persist.CheckIfIdsChanged(fetchIDs)
	if err != nil {
		fmt.Println("Error checking if IDs changed:", err)
	}

	params := config.NewParams(client, corporations, alliances, characters, year, esiData, idChanged, newIDs)

	// Fetch missing data if necessary
	newData, err := GetMissingData(params, dataAvailability)
	if err != nil {
		return nil, err
	}

	// Combine new data with any previously fetched data for a complete year-to-date view
	for month := int(startDate.Month()); month <= int(endDate.Month()); month++ {
		if dataAvailability[month] {
			// Load existing data from file
			fileName := persist.GenerateZkillFileName(year, month)
			monthData, err := persist.ReadKillMailDataFromFile(fileName)
			if err != nil {
				fmt.Println("Error loading data from file:", err)
				continue
			}
			newData = AggregateKillMailDumps(newData, monthData)
		}
	}

	chartData := &model.ChartData{
		KillMails: newData.KillMails,
		ESIData:   *esiData,
	}

	// Refresh ESI data if necessary
	if esiRefresh {
		RefreshCharacter(chartData, client)
	}

	err = persist.SaveEsiDataToFile(persist.GenerateEsiDataFileName(), &chartData.ESIData)
	if err != nil {
		fmt.Println(fmt.Println("Error saving esi data:", err))
	}

	err = persist.SaveIdsToFile(fetchIDs)
	if err != nil {
		fmt.Println(fmt.Println("Error saving ids data:", err))
	}

	fetchTotalTime := time.Since(fetchStart)
	fmt.Println(fmt.Sprintf("Data fetching complete in %f seconds", fetchTotalTime.Seconds()))
	return chartData, nil
}

func GetMissingData(params config.Params, dataAvailability map[int]bool) (*model.KillMailData, error) {
	aggregatedData := &model.KillMailData{
		KillMails: []model.DetailedKillMail{},
	}

	for month, available := range dataAvailability {
		var err error

		if available && !params.ChangedIDs {
			continue
		}

		if params.ChangedIDs {
			// do full pull if IDs have changed
			available = false
		}

		tempParams := params

		newData, err := GetDataForMonth(tempParams, month)
		if err != nil {
			return nil, err
		}

		fileName := persist.GenerateZkillFileName(params.Year, month)
		if err = persist.SaveKillMailsToFile(fileName, newData); err != nil {
			return nil, fmt.Errorf("failed to save fetched data: %w", err)
		}
		aggregatedData = AggregateKillMailDumps(aggregatedData, newData)
	}

	return aggregatedData, nil
}

// GetDataForMonth fetches and processes data for a specific month when it's not already downloaded.
func GetDataForMonth(params config.Params, month int) (*model.KillMailData, error) {
	aggregatedMonthData := &model.KillMailData{
		KillMails: []model.DetailedKillMail{},
	}

	// Create a map to keep track of the kill mail IDs that have been added
	killMailIDs := make(map[int]bool)

	// Fetch data for all entity types and aggregate
	entityGroups := map[string][]int{
		config.EntityTypeCorporation: params.Corporations,
		config.EntityTypeAlliance:    params.Alliances,
		config.EntityTypeCharacter:   params.Characters,
	}

	fmt.Println(fmt.Sprintf("Fetching data for %04d-%02d...", params.Year, month))

	// Assume FetchEntityPageData is a generic function to fetch data for corporations, alliances, or characters
	fetchEntityPageData := func(entityType string, entityID int, page int) ([]model.KillMail, error) {
		var err error
		var killMails []model.KillMail
		switch entityType {
		case config.EntityTypeCorporation:
			killMails, err = zkill.GetCorporatePageData(config.ZkillURL, params.Client, entityID, page, params.Year, month)
		case config.EntityTypeAlliance:
			killMails, err = zkill.GetAlliancePageData(config.ZkillURL, params.Client, entityID, page, params.Year, month)
		case config.EntityTypeCharacter:
			killMails, err = zkill.GetCharacterPageData(config.ZkillURL, params.Client, entityID, page, params.Year, month)
		}
		fmt.Println(fmt.Sprintf("FetchEntityPageData %d killmails for %s %d, page %d", len(killMails), entityType, entityID, page))
		return killMails, err
	}

	rawKillMails := GetRawKillMails(entityGroups, fetchEntityPageData)

	AggregateKillMails(params, rawKillMails, killMailIDs, aggregatedMonthData)

	GetVictimKillMails(params, month, aggregatedMonthData, killMailIDs)

	fmt.Println(fmt.Sprintf("For month %04d-%02d, fetched %d killmails", params.Year, month, len(aggregatedMonthData.KillMails)))
	fmt.Println(fmt.Sprintf("Currently %d characters, %d corporations, and %d alliances", len(params.EsiData.CharacterInfos), len(params.EsiData.CorporationInfos), len(params.EsiData.AllianceInfos)))

	if len(aggregatedMonthData.KillMails) == 0 {
		fmt.Println(fmt.Sprintf("No killmails found for %04d-%02d", params.Year, month))
		fmt.Println("Returning empty data")
	}

	return aggregatedMonthData, nil
}
