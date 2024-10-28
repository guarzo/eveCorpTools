package visuals

import (
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

// CharacterKillData holds the data for character kill counts
type CharacterKillData struct {
	CharacterID int
	KillCount   int
	Name        string
}

func RenderTopCharacters(chartData *model.ChartData) *charts.Bar {
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

	// Convert the map to a slice of CharacterKillData and sort by kill count
	var sortedData []CharacterKillData
	for _, data := range characterKills {
		sortedData = append(sortedData, data)
	}
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].KillCount > sortedData[j].KillCount
	})

	//// Limit to the top 20 characters
	if len(sortedData) > 20 {
		sortedData = sortedData[:20]
	}

	// Prepare data for the chart
	var characterNames []string
	var counts []opts.BarData
	for i, data := range sortedData {
		characterNames = append(characterNames, data.Name)
		counts = append(counts, opts.BarData{
			Value: data.KillCount,
			ItemStyle: &opts.ItemStyle{
				Color: colors[i%len(colors)],
			},
		})
	}

	// Create a new bar chart instance
	bar := newBarChart("           Kill Count", false)
	bar.SetXAxis(characterNames).
		AddSeries("", counts)
	return bar
}
