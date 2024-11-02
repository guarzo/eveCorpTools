package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

type CharacterData struct {
	Name       string `json:"Name"`
	FinalBlows int    `json:"FinalBlows"`
	DamageDone int    `json:"DamageDone"`
}

func GetDamageAndFinalBlows(chartData *model.ChartData) []CharacterData {
	characterStats := make(map[int]*CharacterData)

	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterID := attacker.CharacterID
			if characterID == 0 {
				continue
			}

			// Check if the character is one of ours
			if !config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				continue
			}

			// Get character info
			characterInfo := chartData.CharacterInfos[characterID]

			// Initialize character data if not exists
			data, exists := characterStats[characterID]
			if !exists {
				data = &CharacterData{
					Name: characterInfo.Name,
				}
				characterStats[characterID] = data
				logger.Infof("Character %s included", characterInfo.Name)
			}

			// Accumulate damage done
			data.DamageDone += attacker.DamageDone

			// Check for final blow
			if attacker.FinalBlow {
				data.FinalBlows++
			}
		}
	}

	// Convert map to slice
	var result []CharacterData
	for _, data := range characterStats {
		result = append(result, *data)
	}

	// Sort by damage done descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].DamageDone > result[j].DamageDone
	})

	// Limit to top 20 characters (if desired)
	if len(result) > 20 {
		result = result[:20]
	}

	return result
}
