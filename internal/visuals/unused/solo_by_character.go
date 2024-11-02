package unused

import (
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/visuals"
)

func GetSoloKills(chartData *model.ChartData) *charts.Bar {
	// Initialize a map to count final blows by each attacking character
	characterKills := make(map[string]int)

	// Populate the characterFinalBlows map using attackers from detailed killmails
	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterInfo, exists := chartData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			characterName := characterInfo.Name

			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				// Only increment the kill count if the kill was a solo kill
				if km.ZKB.Solo {
					characterKills[characterName]++
				}
			}
		}
	}

	// Convert the map to a slice of CharacterKillData and sort by final blow count
	var characterData []visuals.CharacterKillData
	for character, solos := range characterKills {
		characterData = append(characterData, visuals.CharacterKillData{
			Name:      character,
			KillCount: solos,
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
				Color: visuals.colors[i%len(visuals.colors)],
			},
		})
	}

	// Create a new bar chart instance
	bar := visuals.newBarChart("Solo Kills", false)
	bar.SetXAxis(sortedCharacters).
		AddSeries("Solo Kills", counts)
	return bar
}
