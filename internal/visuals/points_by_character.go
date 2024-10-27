package visuals

import (
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/model"
)

func RenderPointsPerCharacter(chartData *model.ChartData) *charts.Bar {
	// Initialize a map to count points by each attacking character
	characterPoints := make(map[string]int)

	// Populate the characterPoints map using attackers from detailed killmails
	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterInfo, exists := chartData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			characterName := characterInfo.Name

			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				characterPoints[characterName] += km.ZKB.Points
			}
		}
	}

	// Convert the map to a slice of CharacterKillData and sort by points
	var characterData []CharacterKillData
	for character, points := range characterPoints {
		characterData = append(characterData, CharacterKillData{
			Name:      character,
			KillCount: points,
		})
	}
	sort.Slice(characterData, func(i, j int) bool {
		return characterData[i].KillCount > characterData[j].KillCount
	})

	// Replace the sorted list of character names with the names from the sorted CharacterKillData slice
	sortedCharacters := make([]string, len(characterData))
	for i, data := range characterData {
		sortedCharacters[i] = data.Name
	}

	// Prepare data for the chart
	var counts []opts.BarData
	for i, data := range characterData {
		counts = append(counts, opts.BarData{Value: data.KillCount,
			ItemStyle: &opts.ItemStyle{
				Color: colors[i%len(colors)],
			},
		})
	}

	// Create a new bar chart instance
	bar := newBarChart("Points", false)
	bar.SetXAxis(sortedCharacters).
		AddSeries("Points", counts)
	return bar
}
