package unused

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/visuals"
)

func PrepareKillCountChartData(chartData *model.ChartData) (visuals.ChartJSData, error) {
	characterKills := make(map[int]visuals.CharacterPerformanceData)

	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterID := attacker.CharacterID
			if characterID == 0 {
				continue
			}

			if !config.DisplayCharacter(characterID, attacker.CorporationID, attacker.AllianceID) {
				continue
			}

			if data, found := characterKills[characterID]; found {
				data.KillCount++
				characterKills[characterID] = data
			} else {
				if character, exists := chartData.CharacterInfos[characterID]; exists {
					characterKills[characterID] = visuals.CharacterPerformanceData{
						CharacterID: characterID,
						KillCount:   1,
						Name:        character.Name,
					}
				}
			}
		}
	}

	var sortedData []visuals.CharacterPerformanceData
	for _, data := range characterKills {
		sortedData = append(sortedData, data)
	}
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].KillCount > sortedData[j].KillCount
	})

	// Limit to top 20 characters
	if len(sortedData) > 20 {
		sortedData = sortedData[:20]
	}

	// Prepare data for Chart.js
	var labels []string
	var dataValues []int
	var backgroundColors []string

	for i, d := range sortedData {
		labels = append(labels, d.Name)
		dataValues = append(dataValues, d.KillCount)
		backgroundColors = append(backgroundColors, visuals.colors[i%len(visuals.colors)])
	}

	chartJSData := visuals.ChartJSData{
		Labels: labels,
		Datasets: []visuals.ChartJSDataset{
			{
				Label:           "Kill Count",
				Data:            dataValues,
				BackgroundColor: backgroundColors,
			},
		},
	}

	return chartJSData, nil
}
