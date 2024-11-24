package unused

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

type KillLossRatioData struct {
	CharacterName string
	Kills         int
	Losses        int
	Ratio         float64
}

func GetKillLossRatioData(chartData *model.ChartData) []KillLossRatioData {
	characterStats := make(map[string]*KillLossRatioData)

	for _, km := range chartData.KillMails {
		// Process attackers
		for _, attacker := range km.EsiKillMail.Attackers {
			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				characterInfo := chartData.CharacterInfos[attacker.CharacterID]
				characterName := characterInfo.Name

				data, exists := characterStats[characterName]
				if !exists {
					data = &KillLossRatioData{CharacterName: characterName}
					characterStats[characterName] = data
				}
				data.Kills++
			}
		}

		// Process victim
		victim := km.EsiKillMail.Victim
		if config.TrackedCharacterID(victim.CharacterID) || config.TrackedCorporationID(victim.CorporationID) {
			characterInfo := chartData.CharacterInfos[victim.CharacterID]
			characterName := characterInfo.Name

			data, exists := characterStats[characterName]
			if !exists {
				data = &KillLossRatioData{CharacterName: characterName}
				characterStats[characterName] = data
			}
			data.Losses++
		}
	}

	// Calculate ratios and convert to slice
	var ratioData []KillLossRatioData
	for _, data := range characterStats {
		if data.Losses > 0 {
			data.Ratio = float64(data.Kills) / float64(data.Losses)
		} else {
			data.Ratio = float64(data.Kills)
		}
		ratioData = append(ratioData, *data)
	}

	// Sort by ratio
	sort.Slice(ratioData, func(i, j int) bool {
		return ratioData[i].Ratio > ratioData[j].Ratio
	})

	return ratioData
}
