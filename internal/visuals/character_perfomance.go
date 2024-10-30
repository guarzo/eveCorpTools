package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/model"
)

func GetCharacterPerformance(chartData *model.ChartData) []CharacterKillData {
	characterStats := make(map[int]*CharacterKillData)

	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterID := attacker.CharacterID
			if characterID == 0 {
				continue
			}

			if !isOurCharacter(characterID) {
				continue
			}

			// Get character info
			characterInfo := chartData.CharacterInfos[characterID]
			//if characterInfo == nil {
			//	continue
			//}

			data, exists := characterStats[characterID]
			if !exists {
				data = &CharacterKillData{
					Name: characterInfo.Name,
				}
				characterStats[characterID] = data
			}
			data.KillCount++
			data.Points += km.ZKB.Points

			// Check for solo kill
			if len(km.EsiKillMail.Attackers) == 1 {
				data.SoloKills++
			}
		}
	}

	// Convert map to slice
	var result []CharacterKillData
	for _, data := range characterStats {
		result = append(result, *data)
	}

	// Sort by kill count
	sort.Slice(result, func(i, j int) bool {
		return result[i].KillCount > result[j].KillCount
	})

	return result
}
