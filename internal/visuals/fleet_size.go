package visuals

import (
	"sort"
	"time"

	"github.com/guarzo/zkillanalytics/internal/model"
)

type FleetSizeData struct {
	Time      time.Time
	FleetSize int
}

func GetAverageFleetSizeOverTime(chartData *model.ChartData, interval string) []FleetSizeData {
	fleetSizeMap := make(map[time.Time][]int)

	for _, km := range chartData.KillMails {
		timestamp := km.EsiKillMail.KillMailTime
		var bucket time.Time

		switch interval {
		case "daily":
			bucket = timestamp.Truncate(24 * time.Hour)
		case "weekly":
			year, week := timestamp.ISOWeek()
			bucket = time.Date(year, 0, (week-1)*7+1, 0, 0, 0, 0, time.UTC)
		}

		fleetSizeMap[bucket] = append(fleetSizeMap[bucket], len(km.EsiKillMail.Attackers))
	}

	// Calculate averages
	var fleetSizeData []FleetSizeData
	for t, sizes := range fleetSizeMap {
		total := 0
		for _, size := range sizes {
			total += size
		}
		avgSize := total / len(sizes)
		fleetSizeData = append(fleetSizeData, FleetSizeData{Time: t, FleetSize: avgSize})
	}

	// Sort by time
	sort.Slice(fleetSizeData, func(i, j int) bool {
		return fleetSizeData[i].Time.Before(fleetSizeData[j].Time)
	})

	return fleetSizeData
}
