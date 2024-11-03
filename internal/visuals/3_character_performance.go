package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

// CharacterPerformanceData holds the data for character kill counts
type CharacterPerformanceData struct {
	CharacterID int
	KillCount   int
	Name        string
	Points      int
	SoloKills   int
}

type ChartJSData struct {
	Labels   []string         `json:"labels"`
	Datasets []ChartJSDataset `json:"datasets"`
}

type ChartJSDataset struct {
	Label           string   `json:"label"`
	Data            []int    `json:"data"`
	BackgroundColor []string `json:"backgroundColor"`
}

func GetCharacterPerformance(chartData *model.ChartData) []CharacterPerformanceData {
	characterStats := make(map[int]*CharacterPerformanceData)

	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterID := attacker.CharacterID
			if characterID == 0 {
				continue
			}

			if !config.DisplayCharacter(characterID, attacker.CorporationID, attacker.AllianceID) {
				continue
			}

			// Get character info
			characterInfo := chartData.CharacterInfos[characterID]
			//if characterInfo == nil {
			//	continue
			//}

			data, exists := characterStats[characterID]
			if !exists {
				data = &CharacterPerformanceData{
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
	var result []CharacterPerformanceData
	for _, data := range characterStats {
		result = append(result, *data)
	}

	// Sort by kill count
	sort.Slice(result, func(i, j int) bool {
		return result[i].KillCount > result[j].KillCount
	})

	return result
}
