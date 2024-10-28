package visuals

import (
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
)

// CorporationKillMailData holds the data for corporation kill counts
type CorporationKillMailData struct {
	CorporationID int
	KillCount     int
	Name          string
}

func RenderVictims(chartData *model.ChartData) *charts.Bar {
	// Initialize a map to count killmails by victim corporation
	corpKillMails := make(map[int]CorporationKillMailData)

	// Populate the kill count map using victims from detailed killmails
	for _, km := range chartData.KillMails {
		victimCorpID := km.EsiKillMail.Victim.CorporationID
		if persist.Contains(config.CorporationIDs, victimCorpID) {
			continue
		}
		corpInfo, exists := chartData.CorporationInfos[victimCorpID]
		if !exists || persist.Contains(config.AllianceIDs, corpInfo.AllianceID) {
			continue
		}

		if data, found := corpKillMails[victimCorpID]; found {
			data.KillCount++
			corpKillMails[victimCorpID] = data
		} else {
			corpKillMails[victimCorpID] = CorporationKillMailData{
				CorporationID: victimCorpID,
				KillCount:     1,
				Name:          chartData.CorporationInfos[victimCorpID].Ticker,
			}
		}
	}

	// Convert the map to a slice of CorporationKillMailData and sort by kill count
	var sortedData []CorporationKillMailData
	for _, data := range corpKillMails {
		sortedData = append(sortedData, data)
	}
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].KillCount > sortedData[j].KillCount
	})

	// Limit to the top 10 corporations
	if len(sortedData) > 15 {
		sortedData = sortedData[:15]
	}

	// Prepare data for the chart
	var corpNames []string
	var counts []opts.BarData
	for i, data := range sortedData {
		truncatedName := truncateString(data.Name, 15) // Truncate the name to a maximum of 10 characters
		corpNames = append(corpNames, truncatedName)
		counts = append(counts, opts.BarData{
			Value: data.KillCount,
			ItemStyle: &opts.ItemStyle{
				Color: colors[i%len(colors)], // Assign a color from the list
			},
		})
	}

	// Create a new bar chart instance
	bar := newBarChart("Victims by Corporation", false)
	bar.SetXAxis(corpNames).
		AddSeries("Killmails", counts)
	return bar
}
