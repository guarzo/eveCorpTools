package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

// CharacterKillData holds the data for character kill counts
type CharacterKillData struct {
	CharacterID int
	KillCount   int
	Name        string
}

func PrepareKillCountChartData(chartData *model.ChartData) (ChartJSData, error) {
	// Initialize a map to count kills by character
	characterKills := make(map[int]CharacterKillData)

	// Populate the kill count map using all attackers from detailed killmails
	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			if !config.ExcludeCharacterID(attacker.CharacterID) {
				if data, found := characterKills[attacker.CharacterID]; found {
					data.KillCount++
					characterKills[attacker.CharacterID] = data
				} else {
					if character, exists := chartData.CharacterInfos[attacker.CharacterID]; exists {
						characterKills[attacker.CharacterID] = CharacterKillData{
							CharacterID: attacker.CharacterID,
							KillCount:   1,
							Name:        character.Name,
						}
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
	var data []int
	var backgroundColors []string

	for i, d := range sortedData {
		labels = append(labels, d.Name)
		data = append(data, d.KillCount)
		backgroundColors = append(backgroundColors, colors[i%len(colors)])
	}

	chartJSData := ChartJSData{
		Labels: labels,
		Datasets: []ChartJSDataset{
			{
				Label:           "Kill Count",
				Data:            data,
				BackgroundColor: backgroundColors,
			},
		},
	}

	return chartJSData, nil
}
