package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gambtho/zkillanalytics/internal/model"
)

// fetchPageData fetches a single page of data given a URL
func fetchPageData(client *http.Client, url string) ([]model.KillMail, error) {
	result, err := retryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data from URL: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var killMails []model.KillMail
		if err := json.Unmarshal(body, &killMails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}

		return killMails, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]model.KillMail), nil
}

// fetchKillMails retrieves killmails using a common page fetcher
func fetchKillMails(client *http.Client, baseURL, entity string, entityID, pageNumber, year, month int) ([]model.KillMail, error) {
	url := fmt.Sprintf("%s/api/kills/%s/%d/year/%d/month/%d/page/%d/", baseURL, entity, entityID, year, month, pageNumber)
	return fetchPageData(client, url)
}

// fetchKillMails retrieves killmails using a common page fetcher
func fetchLossKillMails(client *http.Client, baseURL, entity string, entityID, pageNumber, year, month int) ([]model.KillMail, error) {
	url := fmt.Sprintf("%s/api/losses/%s/%d/year/%d/month/%d/page/%d/", baseURL, entity, entityID, year, month, pageNumber)
	return fetchPageData(client, url)
}

// fetchCorporationPageData fetches a single page of killmails for a specific corporation
func fetchCorporationPageData(baseURL string, client *http.Client, corporationID, pageNumber, year, month int) ([]model.KillMail, error) {
	return fetchKillMails(client, baseURL, "corporationID", corporationID, pageNumber, year, month)
}

// fetchAlliancePageData fetches a single page of killmails for a specific alliance
func fetchAlliancePageData(baseURL string, client *http.Client, allianceID, pageNumber, year, month int) ([]model.KillMail, error) {
	return fetchKillMails(client, baseURL, "allianceID", allianceID, pageNumber, year, month)
}

// fetchCharacterPageData fetches a single page of killmails for a specific character
func fetchCharacterPageData(baseURL string, client *http.Client, characterID, pageNumber, year, month int) ([]model.KillMail, error) {
	return fetchKillMails(client, baseURL, "characterID", characterID, pageNumber, year, month)
}

// fetchLossageData fetches a single page of killmails for a specific character
func fetchLossPageData(baseURL string, client *http.Client, characterID, pageNumber, year, month int) ([]model.KillMail, error) {
	return fetchLossKillMails(client, baseURL, "characterID", characterID, pageNumber, year, month)
}

func fetchVictimKillMails(params FetchParams, month int, aggregatedMonthData *model.KillMailData, killMailIDs map[int]bool) {
	trackedCharacters := GetTrackedCharacters(aggregatedMonthData.KillMails, params.EsiData)
	victimKillMails := processTrackedLosses(trackedCharacters, month, params)

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

		_, err := fetchFullKillMail(params.Client, km, aggregatedMonthData)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error processing loss kill mail: %s", err))
			continue
		}
		killMailIDs[int(km.KillMailID)] = true
	}
}

func processTrackedLosses(trackedCharacters []int, month int, params FetchParams) []model.KillMail {
	var victimKillMails []model.KillMail

	// Loop through the trackedCharacters map
	for _, characterID := range trackedCharacters {
		page := 1
		for {
			// Fetch the character page data
			killMails, err := fetchLossPageData(zKillURL, params.Client, characterID, page, params.Year, month)
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

func fetchRawKillMails(entityGroups map[string][]int, fetchEntityPageData func(entityType string, entityID int, page int) ([]model.KillMail, error)) []model.KillMail {
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
			fmt.Println(fmt.Sprintf("FetchRawKillMails %d killmails for %s %d %d", len(rawKillMails), entityName, entityID, page))
		}
	}
	return rawKillMails
}
