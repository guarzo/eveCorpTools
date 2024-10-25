package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gambtho/zkillanalytics/internal/model"
)

// fetchEsiKillMail retrieves detailed killmail information using the killmail ID and hash
func fetchEsiKillMail(baseURL string, client *http.Client, killMailID int64, hash string) (*model.EsiKillMail, error) {
	url := fmt.Sprintf("%s/%d/%s/?datasource=tranquility", baseURL, killMailID, hash)

	result, err := retryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch full killmail data: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var fullKillMail model.EsiKillMail
		if err := json.Unmarshal(body, &fullKillMail); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}

		return &fullKillMail, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.EsiKillMail), nil
}

// FetchCharacterInfo retrieves character information using the character ID
func FetchCharacterInfo(client *http.Client, characterID int) (*model.Character, error) {
	url := fmt.Sprintf("%s/characters/%d/?datasource=tranquility", baseESIURL, characterID)

	result, err := retryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch character data: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var character model.Character
		if err := json.Unmarshal(body, &character); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}
		return &character, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Character), nil

}

// FetchAllianceInfo retrieves alliance information using the alliance ID
func FetchAllianceInfo(client *http.Client, allianceID int) (*model.Alliance, error) {
	url := fmt.Sprintf("%s/alliances/%d/?datasource=tranquility", baseESIURL, allianceID)

	result, err := retryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch alliance data: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var alliance model.Alliance
		if err := json.Unmarshal(body, &alliance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}

		return &alliance, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Alliance), nil
}

func FetchCorporationInfo(client *http.Client, corporationID int) (*model.Corporation, error) {
	url := fmt.Sprintf("%s/corporations/%d/?datasource=tranquility", baseESIURL, corporationID)

	result, err := retryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch corporation data: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var corporation model.Corporation
		if err := json.Unmarshal(body, &corporation); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}

		return &corporation, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Corporation), nil
}

func fetchKillMailDetails(params FetchParams, rawKillMails []model.KillMail, killMailIDs map[int]bool, aggregatedMonthData *model.KillMailData) {
	for index, km := range rawKillMails {
		// Check if the kill mail ID is already in the map
		if _, exists := killMailIDs[int(km.KillMailID)]; exists {
			continue
		}

		// If it's not, process the kill mail and add its ID to the map
		if index%100 == 0 {
			fmt.Println(fmt.Sprintf("Processing killmail %d...%d of %d", km.KillMailID, index, len(rawKillMails)))
		}
		err := processFullKillMail(params.Client, km, aggregatedMonthData, params.EsiData)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error processing kill mail: %s", err))
			continue
		}
		killMailIDs[int(km.KillMailID)] = true
	}
}

func fetchAllESI(client *http.Client, fullKillMail *model.EsiKillMail, esiData *model.ESIData) error {
	if fullKillMail.KillMailID == 0 {
		return nil
	}

	// Add victim to the character map
	if _, exists := esiData.CharacterInfos[fullKillMail.Victim.CharacterID]; !exists {
		if fullKillMail.Victim.CharacterID != 0 {
			victimInfo, err := FetchCharacterInfo(client, fullKillMail.Victim.CharacterID)
			if err != nil {
				return fmt.Errorf("failed to fetch character info for victim %d on km %d: %s", fullKillMail.Victim.CharacterID, fullKillMail.KillMailID, err)
			}
			esiData.CharacterInfos[fullKillMail.Victim.CharacterID] = *victimInfo
		}
	}

	// Add attackers to the character map
	for _, attacker := range fullKillMail.Attackers {
		if attacker.CharacterID != 0 {
			if _, exists := esiData.CharacterInfos[attacker.CharacterID]; !exists {
				characterInfo, err := FetchCharacterInfo(client, attacker.CharacterID)
				if err != nil {
					return fmt.Errorf("failed to fetch character info for attacker %d on km %d: %s", attacker.CharacterID, fullKillMail.KillMailID, err)
				}
				esiData.CharacterInfos[attacker.CharacterID] = *characterInfo
			}
		}
	}

	// Add corporations and alliances to their respective maps
	for _, character := range esiData.CharacterInfos {
		if character.CorporationID != 0 {
			if _, exists := esiData.CorporationInfos[character.CorporationID]; !exists {
				corporationInfo, err := FetchCorporationInfo(client, character.CorporationID)
				if err != nil {
					return fmt.Errorf("failed to fetch corporation info for corporation %d: %s", character.CorporationID, err)
				}
				esiData.CorporationInfos[character.CorporationID] = *corporationInfo
			}
		}
	}

	for _, corporation := range esiData.CorporationInfos {
		if corporation.AllianceID != 0 {
			if _, exists := esiData.AllianceInfos[corporation.AllianceID]; !exists {
				allianceInfo, err := FetchAllianceInfo(client, corporation.AllianceID)
				if err != nil {
					return fmt.Errorf("failed to fetch alliance info for alliance %d: %s", corporation.AllianceID, err)
				}
				esiData.AllianceInfos[corporation.AllianceID] = *allianceInfo
			}
		}
	}
	return nil
}
