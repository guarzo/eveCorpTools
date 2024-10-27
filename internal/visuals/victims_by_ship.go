package visuals

import (
	"fmt"
	"sort"

	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/service"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/gambtho/zkillanalytics/internal/model"
)

// ShipKillData holds the data for ship kill counts
type ShipKillData struct {
	ShipTypeID int
	KillCount  int
	Name       string
}

func RenderTopShipsKilled(orchestrator *service.OrchestrateService, chartData *model.ChartData) *charts.Bar {
	// Initialize a map to count killmails by ship type
	shipKillCounts := make(map[int]ShipKillData)

	if trackedCharacters == nil || len(trackedCharacters) == 0 {
		fmt.Print(fmt.Sprintf("No tracked characters found, fetching from %d killmails", len(chartData.KillMails)))
		trackedCharacters = orchestrator.GetTrackedCharactersFromKillMails(chartData.KillMails, &chartData.ESIData)
	}

	// Populate the kill count map using victims' ships from detailed killmails
	for _, km := range chartData.KillMails {
		victim := km.EsiKillMail.Victim
		if persist.Contains(trackedCharacters, victim.CharacterID) {
			continue
		}

		shipTypeID := km.EsiKillMail.Victim.ShipTypeID
		shipName := orchestrator.LookupType(shipTypeID) // Fetch the ship name

		if shipName == "" || shipName == "Capsule" || shipName == "#System" || shipName == "Mobile Tractor Unit" {
			continue
		}

		if data, found := shipKillCounts[shipTypeID]; found {
			data.KillCount++
			shipKillCounts[shipTypeID] = data
		} else {
			shipKillCounts[shipTypeID] = ShipKillData{
				ShipTypeID: shipTypeID,
				KillCount:  1,
				Name:       shipName,
			}
		}
	}

	// Convert the map to a slice of ShipKillData and sort by kill count
	var sortedData []ShipKillData
	for _, data := range shipKillCounts {
		sortedData = append(sortedData, data)
	}
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].KillCount > sortedData[j].KillCount
	})

	// Limit to the top 10 ships
	if len(sortedData) > 20 {
		sortedData = sortedData[:20]
	}

	// Prepare data for the chart
	var shipNames []string
	var counts []opts.BarData
	for i, data := range sortedData {
		truncatedName := truncateString(data.Name, 15) // Truncate the name to a maximum of 15 characters
		shipNames = append(shipNames, truncatedName)
		counts = append(counts, opts.BarData{
			Value: data.KillCount,
			ItemStyle: &opts.ItemStyle{
				Color: colors[i%len(colors)], // Assign a color from the list
			},
		})
	}

	// Create a new bar chart instance
	bar := newBarChart("Top Ships Killed", false)
	bar.SetXAxis(shipNames).
		AddSeries("Killmails", counts)
	return bar
}
