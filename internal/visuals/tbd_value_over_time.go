package visuals

import (
	"sort"
	"time"

	"github.com/guarzo/zkillanalytics/internal/model"
)

type ValueOverTimeData struct {
	Time  time.Time
	Value float64
}

func GetValueOverTimeData(chartData *model.ChartData, interval string) []ValueOverTimeData {
	valueMap := make(map[time.Time]float64)

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

		valueMap[bucket] += km.ZKB.TotalValue
	}

	// Convert map to slice and sort
	var valueSeries []ValueOverTimeData
	for t, value := range valueMap {
		valueSeries = append(valueSeries, ValueOverTimeData{Time: t, Value: value})
	}
	sort.Slice(valueSeries, func(i, j int) bool {
		return valueSeries[i].Time.Before(valueSeries[j].Time)
	})

	return valueSeries
}
