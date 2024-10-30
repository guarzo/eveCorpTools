package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

type CharacterPerformanceData struct {
	Name      string
	KillCount int
	SoloKills int
	Points    int
}

func GetCharacterPerformanceData(chartData *model.ChartData) []CharacterPerformanceData {
	characterDataMap := make(map[string]*CharacterPerformanceData)

	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				characterInfo := chartData.CharacterInfos[attacker.CharacterID]
				characterName := characterInfo.Name

				data, exists := characterDataMap[characterName]
				if !exists {
					data = &CharacterPerformanceData{Name: characterName}
					characterDataMap[characterName] = data
				}

				data.KillCount++
				data.Points += km.ZKB.Points

				if km.ZKB.Solo {
					data.SoloKills++
				}
			}
		}
	}

	var performanceDataSlice []CharacterPerformanceData
	for _, data := range characterDataMap {
		performanceDataSlice = append(performanceDataSlice, *data)
	}

	// Sort by kill count
	sort.Slice(performanceDataSlice, func(i, j int) bool {
		return performanceDataSlice[i].KillCount > performanceDataSlice[j].KillCount
	})

	return performanceDataSlice
}
