package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

type KillLossAndISKEfficiencyData struct {
	CharacterName string
	Kills         int
	Losses        int
	Ratio         float64
	ISKDestroyed  float64
	ISKLost       float64
	Efficiency    float64
}

func GetKillLossAndISKEfficiencyData(chartData *model.ChartData) []KillLossAndISKEfficiencyData {
	characterStats := make(map[string]*KillLossAndISKEfficiencyData)

	for _, km := range chartData.KillMails {
		// Process attackers
		for _, attacker := range km.EsiKillMail.Attackers {
			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				characterInfo := chartData.CharacterInfos[attacker.CharacterID]
				characterName := characterInfo.Name

				data, exists := characterStats[characterName]
				if !exists {
					data = &KillLossAndISKEfficiencyData{CharacterName: characterName}
					characterStats[characterName] = data
				}
				data.Kills++
				data.ISKDestroyed += km.ZKB.TotalValue
			}
		}

		// Process victim
		victim := km.EsiKillMail.Victim
		if config.TrackedCharacterID(victim.CharacterID) || config.TrackedCorporationID(victim.CorporationID) {
			characterInfo := chartData.CharacterInfos[victim.CharacterID]
			characterName := characterInfo.Name

			data, exists := characterStats[characterName]
			if !exists {
				data = &KillLossAndISKEfficiencyData{CharacterName: characterName}
				characterStats[characterName] = data
			}
			data.Losses++
			data.ISKLost += km.ZKB.TotalValue
		}
	}

	// Calculate Ratio and Efficiency
	var statsData []KillLossAndISKEfficiencyData
	for _, data := range characterStats {
		if data.Losses > 0 {
			data.Ratio = float64(data.Kills) / float64(data.Losses)
		} else if data.Kills > 0 {
			data.Ratio = float64(data.Kills)
		} else {
			data.Ratio = 0
		}

		totalISK := data.ISKDestroyed + data.ISKLost
		if totalISK > 0 {
			data.Efficiency = (data.ISKDestroyed / totalISK) * 100
		} else {
			data.Efficiency = 0
		}
		statsData = append(statsData, *data)
	}

	// Sort by Ratio
	sort.Slice(statsData, func(i, j int) bool {
		return statsData[i].Ratio > statsData[j].Ratio
	})

	return statsData
}
