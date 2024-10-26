package service

import (
	"fmt"
	"net/http"

	"github.com/gambtho/zkillanalytics/internal/api/esi"
	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

func AggregateEsi(client *http.Client, fullKillMail *model.EsiKillMail, esiData *model.ESIData) error {
	if fullKillMail.KillMailID == 0 {
		return nil
	}

	// Add victim to the character map
	if _, exists := esiData.CharacterInfos[fullKillMail.Victim.CharacterID]; !exists {
		if fullKillMail.Victim.CharacterID != 0 {
			victimInfo, err := esi.FetchCharacterInfo(client, fullKillMail.Victim.CharacterID)
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
				characterInfo, err := esi.FetchCharacterInfo(client, attacker.CharacterID)
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
				corporationInfo, err := esi.GetCorporationInfo(client, character.CorporationID)
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
				allianceInfo, err := esi.GetAllianceInfo(client, corporation.AllianceID)
				if err != nil {
					return fmt.Errorf("failed to fetch alliance info for alliance %d: %s", corporation.AllianceID, err)
				}
				esiData.AllianceInfos[corporation.AllianceID] = *allianceInfo
			}
		}
	}
	return nil
}

func RefreshCharacter(chartData *model.ChartData, client *http.Client) {
	fmt.Println("Refreshing ESI data...")

	emptyESI := false

	if len(chartData.ESIData.CharacterInfos) == 0 {
		fmt.Println("Empty ESI file provided to refresh")
		emptyESI = true
	}

	newEsiData := &model.ESIData{
		AllianceInfos:    make(map[int]model.Alliance),
		CharacterInfos:   make(map[int]model.Character),
		CorporationInfos: make(map[int]model.Corporation),
	}

	fmt.Println(fmt.Sprintf("Refreshing ESI data for characters... %d killmails to process", len(chartData.KillMails)))
	for index, detailedKillMail := range chartData.KillMails {
		if index%100 == 0 {
			fmt.Println(fmt.Sprintf("Processing killmail %d...%d of %d", detailedKillMail.EsiKillMail.KillMailID, index, len(chartData.KillMails)))
		}
		err := AggregateEsi(client, &detailedKillMail.EsiKillMail, newEsiData)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to fetch ESI data %s", err))
			if !emptyESI {
				fmt.Println("Error fetching ESI data, using existing data")
				return
			}
		}
	}

	chartData.ESIData = *newEsiData
	fmt.Println("ESI data refreshed.")
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
