package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

type ISKEfficiencyData struct {
	CharacterName string
	ISKDestroyed  float64
	ISKLost       float64
	Efficiency    float64
}

func GetISKEfficiencyData(chartData *model.ChartData) []ISKEfficiencyData {
	characterStats := make(map[string]*ISKEfficiencyData)

	for _, km := range chartData.KillMails {
		// Process attackers
		for _, attacker := range km.EsiKillMail.Attackers {
			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				characterInfo := chartData.CharacterInfos[attacker.CharacterID]
				characterName := characterInfo.Name

				data, exists := characterStats[characterName]
				if !exists {
					data = &ISKEfficiencyData{CharacterName: characterName}
					characterStats[characterName] = data
				}
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
				data = &ISKEfficiencyData{CharacterName: characterName}
				characterStats[characterName] = data
			}
			data.ISKLost += km.ZKB.TotalValue
		}
	}

	// Calculate efficiency and convert to slice
	var efficiencyData []ISKEfficiencyData
	for _, data := range characterStats {
		totalISK := data.ISKDestroyed + data.ISKLost
		if totalISK > 0 {
			data.Efficiency = (data.ISKDestroyed / totalISK) * 100
		} else {
			data.Efficiency = 0
		}
		efficiencyData = append(efficiencyData, *data)
	}

	// Sort by efficiency
	sort.Slice(efficiencyData, func(i, j int) bool {
		return efficiencyData[i].Efficiency > efficiencyData[j].Efficiency
	})

	return efficiencyData
}
