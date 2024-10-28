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

func RenderLostShipTypes(orchestrator *service.OrchestrateService, chartData *model.ChartData) *charts.Bar {
	// Initialize a map to count lost ship types
	shipLosses := make(map[string]int)

	if trackedCharacters == nil || len(trackedCharacters) == 0 {
		fmt.Print(fmt.Sprintf("No tracked characters found, fetching from %d killmails", len(chartData.KillMails)))
		trackedCharacters = orchestrator.GetTrackedCharactersFromKillMails(chartData.KillMails, &chartData.ESIData)
	}

	// Populate the shipLosses map using victims from detailed killmails
	for _, km := range chartData.KillMails {
		victim := km.EsiKillMail.Victim
		if !persist.Contains(trackedCharacters, victim.CharacterID) {
			continue
		}

		shipName := orchestrator.LookupType(victim.ShipTypeID)

		if shipName == "" || shipName == "Capsule" || shipName == "#System" || shipName == "Mobile Tractor Unit" {
			continue
		}

		shipLosses[shipName]++
	}

	// Convert the map to a slice of ShipLossData and sort by loss count
	type ShipLossData struct {
		Name      string
		LossCount int
	}
	var shipData []ShipLossData
	for ship, losses := range shipLosses {

		shipData = append(shipData, ShipLossData{
			Name:      ship,
			LossCount: losses,
		})
	}
	sort.Slice(shipData, func(i, j int) bool {
		return shipData[i].LossCount > shipData[j].LossCount
	})

	//// Limit to the top 20 characters
	if len(shipData) > 20 {
		shipData = shipData[:20]
	}

	// Collect the sorted ship names and their loss counts
	var shipNames []string
	var lossCounts []opts.BarData
	for i, data := range shipData {
		shipNames = append(shipNames, data.Name)
		lossCounts = append(lossCounts, opts.BarData{Value: data.LossCount, ItemStyle: &opts.ItemStyle{
			Color: colors[i%len(colors)],
		}})
	}

	// Create a new bar chart instance
	bar := newBarChart("Ship Losses", false)

	bar.SetXAxis(shipNames).
		AddSeries("Losses", lossCounts)

	return bar
}
