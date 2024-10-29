package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

// Data structure to hold character data
type CharacterData struct {
	Name       string
	FinalBlows int
	DamageDone int
}

// Function to process data
func GetDamageOrFinalBlowsData(chartData *model.ChartData, displayType string) []CharacterData {
	// Initialize a map to count final blows and damage done by each character
	characterDataMap := make(map[string]*CharacterData)

	// Process the killmails
	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterInfo, exists := chartData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			characterName := characterInfo.Name

			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				// Get or create the character data
				data, exists := characterDataMap[characterName]
				if !exists {
					data = &CharacterData{Name: characterName}
					characterDataMap[characterName] = data
				}

				// Update data
				if attacker.FinalBlow {
					data.FinalBlows++
				}
				data.DamageDone += attacker.DamageDone
			}
		}
	}

	// Convert the map to a slice
	var characterDataSlice []CharacterData
	for _, data := range characterDataMap {
		characterDataSlice = append(characterDataSlice, *data)
	}

	// Sort the data
	if displayType == "damage" {
		sort.Slice(characterDataSlice, func(i, j int) bool {
			return characterDataSlice[i].DamageDone > characterDataSlice[j].DamageDone
		})
	} else if displayType == "blows" {
		sort.Slice(characterDataSlice, func(i, j int) bool {
			return characterDataSlice[i].FinalBlows > characterDataSlice[j].FinalBlows
		})
	}

	return characterDataSlice
}
