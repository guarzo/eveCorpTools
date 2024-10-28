package visuals

import (
	"fmt"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/service"
)

type CharacterValueData struct {
	Name  string
	Value float64
}

func RenderOurLossesValue(orchestrateService *service.OrchestrateService, chartData *model.ChartData) *charts.Bar {
	// Initialize a map to sum totalValue by each victim character
	characterValues := make(map[string]float64)

	if trackedCharacters == nil || len(trackedCharacters) == 0 {
		fmt.Print(fmt.Sprintf("No tracked characters found, fetching from %d killmails", len(chartData.KillMails)))
		trackedCharacters = orchestrateService.GetTrackedCharactersFromKillMails(chartData.KillMails, &chartData.ESIData)
	}

	// Populate the characterValues map using victims from detailed killmails
	for _, km := range chartData.KillMails {
		victim := km.EsiKillMail.Victim

		if !persist.Contains(trackedCharacters, victim.CharacterID) {
			continue
		}

		characterInfo, exists := chartData.CharacterInfos[victim.CharacterID]
		if !exists {
			continue
		}

		characterName := characterInfo.Name

		characterValues[characterName] += km.ZKB.TotalValue
	}

	// Convert the map to a slice of CharacterValueData and sort by totalValue
	var characterData []CharacterValueData
	for character, value := range characterValues {
		characterData = append(characterData, CharacterValueData{
			Name:  character,
			Value: value,
		})
	}
	sort.Slice(characterData, func(i, j int) bool {
		return characterData[i].Value > characterData[j].Value
	})

	// Replace the sorted list of character names with the names from the sorted CharacterValueData slice
	sortedCharacters := make([]string, len(characterData))
	for i, data := range characterData {
		sortedCharacters[i] = data.Name
	}

	// Prepare data for the chart
	var values []opts.BarData
	for i, data := range characterData {
		values = append(values, opts.BarData{Value: data.Value,
			ItemStyle: &opts.ItemStyle{
				Color: colors[i%len(colors)],
			},
		})
	}

	// Create a new bar chart instance
	bar := newBarChart("Our Lost Isk", false)
	bar.SetXAxis(sortedCharacters).
		AddSeries("Total Value", values)
	return bar
}
