package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/model"
)

// CharacterKillData holds the data for character kill counts
type CharacterKillData struct {
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

func PrepareKillCountChartData(chartData *model.ChartData) (ChartJSData, error) {
	characterKills := make(map[int]CharacterKillData)

	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterID := attacker.CharacterID
			if characterID == 0 {
				continue
			}

			if !isOurCharacter(characterID) {
				continue
			}

			if data, found := characterKills[characterID]; found {
				data.KillCount++
				characterKills[characterID] = data
			} else {
				if character, exists := chartData.CharacterInfos[characterID]; exists {
					characterKills[characterID] = CharacterKillData{
						CharacterID: characterID,
						KillCount:   1,
						Name:        character.Name,
					}
				}
			}
		}
	}

	var sortedData []CharacterKillData
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
		backgroundColors = append(backgroundColors, colors[i%len(colors)])
	}

	chartJSData := ChartJSData{
		Labels: labels,
		Datasets: []ChartJSDataset{
			{
				Label:           "Kill Count",
				Data:            dataValues,
				BackgroundColor: backgroundColors,
			},
		},
	}

	return chartJSData, nil
}
