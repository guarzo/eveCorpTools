package visuals

import (
	"fmt"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/service"
)

func RenderOurLossesCount(orchestrator *service.OrchestrateService, chartData *model.ChartData) *charts.Bar {
	// Initialize a map to count losses by each victim character
	characterLosses := make(map[string]int)

	if trackedCharacters == nil || len(trackedCharacters) == 0 {
		fmt.Print(fmt.Sprintf("No tracked characters found, fetching from %d killmails", len(chartData.KillMails)))
		trackedCharacters = orchestrator.GetTrackedCharactersFromKillMails(chartData.KillMails, &chartData.ESIData)
	}

	// Populate the characterLosses map using victims from detailed killmails
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
		characterLosses[characterName]++
	}

	// Convert the map to a slice of CharacterKillData and sort by losses
	var characterData []CharacterKillData
	for character, losses := range characterLosses {
		characterData = append(characterData, CharacterKillData{
			Name:      character,
			KillCount: losses,
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
	bar := newBarChart("Our Losses", false)
	bar.SetXAxis(sortedCharacters).
		AddSeries("Losses", counts)
	return bar
}
