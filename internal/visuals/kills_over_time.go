package visuals

import (
	"sort"
	"time"

	"github.com/guarzo/zkillanalytics/internal/model"
)

type TimeSeriesData struct {
	Time  time.Time
	Kills int
}

func GetKillActivityOverTime(chartData *model.ChartData, interval string) []TimeSeriesData {
	killCounts := make(map[time.Time]int)

	for _, km := range chartData.KillMails {
		timestamp := km.KillMailTime
		var bucket time.Time

		switch interval {
		case "hourly":
			bucket = timestamp.Truncate(time.Hour)
		case "daily":
			bucket = timestamp.Truncate(24 * time.Hour)
		case "weekly":
			year, week := timestamp.ISOWeek()
			bucket = time.Date(year, 0, (week-1)*7+1, 0, 0, 0, 0, time.UTC)
		}

		killCounts[bucket]++
	}

	// Convert map to slice and sort by time
	var timeSeries []TimeSeriesData
	for t, count := range killCounts {
		timeSeries = append(timeSeries, TimeSeriesData{Time: t, Kills: count})
	}
	sort.Slice(timeSeries, func(i, j int) bool {
		return timeSeries[i].Time.Before(timeSeries[j].Time)
	})

	return timeSeries
}
