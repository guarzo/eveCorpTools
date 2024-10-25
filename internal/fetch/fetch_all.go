package fetch

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

// FetchAllData Main function to orchestrate data fetching based on availability and necessity.
func FetchAllData(client *http.Client, corporations, alliances, characters []int, startString, endString string) (*model.ChartData, error) {
	fmt.Println(fmt.Sprintf("Fetching data for %s to %s...", startString, endString))
	fetchStart := time.Now()

	startDate, _ := time.Parse("2006-01-02", startString)
	endDate, _ := time.Parse("2006-01-02", endString)
	year := startDate.Year()
	esiRefresh := false

	dataAvailability, _ := checkDataAvailability(int(startDate.Month()), int(endDate.Month()), year)

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

	params := NewFetchParams(client, corporations, alliances, characters, year, esiData, idChanged, newIDs)

	// Fetch missing data if necessary
	newData, err := fetchMissingData(params, dataAvailability)
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
			newData = aggregateKillMailDumps(newData, monthData)
		}
	}

	chartData := &model.ChartData{
		KillMails: newData.KillMails,
		ESIData:   *esiData,
	}

	// Refresh ESI data if necessary
	if esiRefresh {
		RefreshEsiData(chartData, client)
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

// checkDataAvailability checks which months within a range have data files already present.
func checkDataAvailability(startMonth, endMonth, year int) (map[int]bool, error) {
	dataAvailability := make(map[int]bool)
	for month := startMonth; month <= endMonth; month++ {
		fileName := persist.GenerateZkillFileName(year, month)
		if fileInfo, err := os.Stat(fileName); err == nil {
			if fileInfo.Size() <= 1*1024 {
				fmt.Printf("File %s is too small (%d bytes). Marking as unavailable.\n", fileName, fileInfo.Size())
				dataAvailability[month] = false
				continue
			}
			fmt.Println(fmt.Sprintf("Data for %04d-%02d already exists.", year, month))
			dataAvailability[month] = true
			if month == int(time.Now().Month()) {
				fileDate := fileInfo.ModTime().Truncate(24 * time.Hour)
				today := time.Now().Truncate(24 * time.Hour)
				if fileDate.Before(today) {
					fmt.Println(fmt.Sprintf("Removing stale month to date file %s...", fileName))
					dataAvailability[month] = false
					err = os.Remove(fileName)
					if err != nil {
						fmt.Println(fmt.Sprintf("Error removing stale month to date file %s: %s", fileName, err))
					}
				} else {
					fmt.Println(fmt.Sprintf("Continuing to use current month data for %s, %s", fileName, fileDate.Format("2006-01-02:15:04:05")))
				}
			}
		} else {
			dataAvailability[month] = false
		}
	}
	return dataAvailability, nil
}

func fetchMissingData(params FetchParams, dataAvailability map[int]bool) (*model.KillMailData, error) {
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

		newData, err := fetchDataForMonth(tempParams, month)
		if err != nil {
			return nil, err
		}

		fileName := persist.GenerateZkillFileName(params.Year, month)
		if err = persist.SaveKillMailsToFile(fileName, newData); err != nil {
			return nil, fmt.Errorf("failed to save fetched data: %w", err)
		}
		aggregatedData = aggregateKillMailDumps(aggregatedData, newData)
	}

	return aggregatedData, nil
}

// aggregateKillMailDumps combines data from multiple KillMailData objects.
func aggregateKillMailDumps(base, addition *model.KillMailData) *model.KillMailData {
	if base == nil {
		return addition
	}
	if addition == nil {
		return base
	}

	base.KillMails = append(base.KillMails, addition.KillMails...)
	return base
}

// processFullKillMail processes a kill mail and updates the aggregated data
func processFullKillMail(client *http.Client, km model.KillMail, aggregatedData *model.KillMailData, esiData *model.ESIData) error {
	fullKillMail, err := fetchFullKillMail(client, km, aggregatedData)
	if err != nil {
		return err
	}

	err = fetchAllESI(client, fullKillMail, esiData)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to fetch ESI data %s", err))
	}

	return nil
}

func fetchFullKillMail(client *http.Client, km model.KillMail, aggregatedData *model.KillMailData) (*model.EsiKillMail, error) {
	fullKillMail, err := fetchEsiKillMail(esiKillURL, client, km.KillMailID, km.ZKB.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch full killmail for ID %d: %s", km.KillMailID, err)
	}
	dKM := model.DetailedKillMail{KillMail: km, EsiKillMail: *fullKillMail}
	aggregatedData.KillMails = append(aggregatedData.KillMails, dKM)
	return fullKillMail, nil
}

func GetTrackedCharacters(fullKillMail []model.DetailedKillMail, esiData *model.ESIData) []int {
	var trackedCharacters []int

	for _, km := range fullKillMail {
		for _, attacker := range km.Attackers {
			if persist.Contains(trackedCharacters, attacker.CharacterID) {
				continue
			}

			_, exists := esiData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			corpInfo, exists := esiData.CorporationInfos[attacker.CorporationID]
			if !exists {
				continue
			}

			allianceID := corpInfo.AllianceID

			if persist.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, allianceID) {
				// fmt.Println(fmt.Sprintf("Adding character %d to tracked characters", attacker.CharacterID))
				trackedCharacters = append(trackedCharacters, attacker.CharacterID)
			}
		}
	}

	// fmt.Println(fmt.Sprintf("Found %d tracked characters", len(trackedCharacters)))
	return trackedCharacters
}

// fetchDataForMonth fetches and processes data for a specific month when it's not already downloaded.
func fetchDataForMonth(params FetchParams, month int) (*model.KillMailData, error) {
	aggregatedMonthData := &model.KillMailData{
		KillMails: []model.DetailedKillMail{},
	}

	// Create a map to keep track of the kill mail IDs that have been added
	killMailIDs := make(map[int]bool)

	// Fetch data for all entity types and aggregate
	entityGroups := map[string][]int{
		entityTypeCorporation: params.Corporations,
		entityTypeAlliance:    params.Alliances,
		entityTypeCharacter:   params.Characters,
	}

	fmt.Println(fmt.Sprintf("Fetching data for %04d-%02d...", params.Year, month))

	// Assume FetchEntityPageData is a generic function to fetch data for corporations, alliances, or characters
	fetchEntityPageData := func(entityType string, entityID int, page int) ([]model.KillMail, error) {
		var err error
		var killMails []model.KillMail
		switch entityType {
		case entityTypeCorporation:
			killMails, err = fetchCorporationPageData(zKillURL, params.Client, entityID, page, params.Year, month)
		case entityTypeAlliance:
			killMails, err = fetchAlliancePageData(zKillURL, params.Client, entityID, page, params.Year, month)
		case entityTypeCharacter:
			killMails, err = fetchCharacterPageData(zKillURL, params.Client, entityID, page, params.Year, month)
		}
		fmt.Println(fmt.Sprintf("FetchEntityPageData %d killmails for %s %d, page %d", len(killMails), entityType, entityID, page))
		return killMails, err
	}

	rawKillMails := fetchRawKillMails(entityGroups, fetchEntityPageData)

	fetchKillMailDetails(params, rawKillMails, killMailIDs, aggregatedMonthData)

	fetchVictimKillMails(params, month, aggregatedMonthData, killMailIDs)

	fmt.Println(fmt.Sprintf("For month %04d-%02d, fetched %d killmails", params.Year, month, len(aggregatedMonthData.KillMails)))
	fmt.Println(fmt.Sprintf("Currently %d characters, %d corporations, and %d alliances", len(params.EsiData.CharacterInfos), len(params.EsiData.CorporationInfos), len(params.EsiData.AllianceInfos)))

	if len(aggregatedMonthData.KillMails) == 0 {
		fmt.Println(fmt.Sprintf("No killmails found for %04d-%02d", params.Year, month))
		fmt.Println("Returning empty data")
	}

	return aggregatedMonthData, nil
}
