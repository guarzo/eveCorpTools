package service

import (
	"fmt"
	"net/http"

	"github.com/gambtho/zkillanalytics/internal/api/esi"
	"github.com/gambtho/zkillanalytics/internal/api/zkill"
	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/model"
)

func GetVictimKillMails(params config.Params, month int, aggregatedMonthData *model.KillMailData, killMailIDs map[int]bool) {
	trackedCharacters := GetTrackedCharacters(aggregatedMonthData.KillMails, params.EsiData)
	victimKillMails := ProcessTrackedKills(trackedCharacters, month, params)

	fmt.Println(fmt.Sprintf("Processing %d loss killmails...", len(victimKillMails)))
	for index, km := range victimKillMails {
		// Check if the kill mail ID is already in the map
		if _, exists := killMailIDs[int(km.KillMailID)]; exists {
			continue
		}

		// If it's not, process the kill mail and add its ID to the map
		if index%100 == 0 {
			fmt.Println(fmt.Sprintf("Processing loss killmail %d...%d of %d", km.KillMailID, index, len(victimKillMails)))
		}

		_, err := GetFullKillMail(params.Client, km, aggregatedMonthData)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error processing loss kill mail: %s", err))
			continue
		}
		killMailIDs[int(km.KillMailID)] = true
	}
}

func ProcessTrackedKills(trackedCharacters []int, month int, params config.Params) []model.KillMail {
	var victimKillMails []model.KillMail

	// Loop through the trackedCharacters map
	for _, characterID := range trackedCharacters {
		page := 1
		for {
			// Fetch the character page data
			killMails, err := zkill.GetLossPageData(config.ZkillURL, params.Client, characterID, page, params.Year, month)
			if err != nil {
				fmt.Printf("Error fetching data for characterID %d: %v\n", characterID, err)
				break
			}

			// Break the loop if the page is empty
			if len(killMails) == 0 {
				break
			}

			// Append the fetched kill mails to victimKillMails
			victimKillMails = append(victimKillMails, killMails...)
			page++
		}
	}
	// fmt.Println(fmt.Sprintf("Fetched %d loss killmails for all tracked characters", len(victimKillMails)))
	return victimKillMails
}

func GetRawKillMails(entityGroups map[string][]int, fetchEntityPageData func(entityType string, entityID int, page int) ([]model.KillMail, error)) []model.KillMail {
	var rawKillMails []model.KillMail

	for entityName, entityValues := range entityGroups {
		for _, entityID := range entityValues {
			page := 1
			for {
				killMails, err := fetchEntityPageData(entityName, entityID, page)
				if err != nil || len(killMails) == 0 {
					break
				}
				rawKillMails = append(rawKillMails, killMails...)
				page++
			}
			fmt.Println(fmt.Sprintf("GetRawKillMails %d killmails for %s %d %d", len(rawKillMails), entityName, entityID, page))
		}
	}
	return rawKillMails
}

func AggregateKillMails(params config.Params, rawKillMails []model.KillMail, killMailIDs map[int]bool, aggregatedMonthData *model.KillMailData) {
	for index, km := range rawKillMails {
		// Check if the kill mail ID is already in the map
		if _, exists := killMailIDs[int(km.KillMailID)]; exists {
			continue
		}

		// If it's not, process the kill mail and add its ID to the map
		if index%100 == 0 {
			fmt.Println(fmt.Sprintf("Processing killmail %d...%d of %d", km.KillMailID, index, len(rawKillMails)))
		}
		err := ProcessFullKillMail(params.Client, km, aggregatedMonthData, params.EsiData)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error processing kill mail: %s", err))
			continue
		}
		killMailIDs[int(km.KillMailID)] = true
	}
}

// AggregateKillMailDumps combines data from multiple KillMailData objects.
func AggregateKillMailDumps(base, addition *model.KillMailData) *model.KillMailData {
	if base == nil {
		return addition
	}
	if addition == nil {
		return base
	}

	base.KillMails = append(base.KillMails, addition.KillMails...)
	return base
}

// ProcessFullKillMail processes a kill mail and updates the aggregated data
func ProcessFullKillMail(client *http.Client, km model.KillMail, aggregatedData *model.KillMailData, esiData *model.ESIData) error {
	fullKillMail, err := GetFullKillMail(client, km, aggregatedData)
	if err != nil {
		return err
	}

	err = AggregateEsi(client, fullKillMail, esiData)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to fetch ESI data %s", err))
	}

	return nil
}

func GetFullKillMail(client *http.Client, km model.KillMail, aggregatedData *model.KillMailData) (*model.EsiKillMail, error) {
	fullKillMail, err := esi.GetEsiKillMail(config.EsiURL, client, km.KillMailID, km.ZKB.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch full killmail for ID %d: %s", km.KillMailID, err)
	}
	dKM := model.DetailedKillMail{KillMail: km, EsiKillMail: *fullKillMail}
	aggregatedData.KillMails = append(aggregatedData.KillMails, dKM)
	return fullKillMail, nil
}
