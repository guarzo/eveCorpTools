// internal/visuals/fleetSizeAndValue.go

package visuals

import (
	"sort"
	"time"

	"github.com/guarzo/zkillanalytics/internal/model"
)

// FleetSizeAndValueData holds the average fleet size and total value for a specific time bucket
type FleetSizeAndValueData struct {
	Time         time.Time `json:"time"`
	AvgFleetSize float64   `json:"avg_fleet_size"`
	TotalValue   float64   `json:"total_value"`
}

// GetFleetSizeAndValueOverTime calculates the average fleet size and total value over specified time intervals
func GetFleetSizeAndValueOverTime(chartData *model.ChartData, interval string) []FleetSizeAndValueData {
	fleetValueMap := make(map[time.Time]struct {
		TotalFleetSize int
		Count          int
		TotalValue     float64
	})

	for _, km := range chartData.KillMails {
		timestamp := km.EsiKillMail.KillMailTime
		var bucket time.Time

		switch interval {
		case "daily":
			bucket = time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
		case "weekly":
			year, week := timestamp.ISOWeek()
			// ISOWeek: week starts on Monday
			// Calculate the date of the Monday of the ISO week
			bucket = getMondayOfISOWeek(year, week)
		default:
			// Default to daily if interval is unrecognized
			bucket = time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
		}

		fleetSize := len(km.EsiKillMail.Attackers)
		value := km.ZKB.TotalValue

		data := fleetValueMap[bucket]
		data.TotalFleetSize += fleetSize
		data.Count++
		data.TotalValue += value
		fleetValueMap[bucket] = data
	}

	// Convert map to slice and calculate averages
	var fleetValueData []FleetSizeAndValueData
	for t, data := range fleetValueMap {
		avgFleetSize := 0.0
		if data.Count > 0 {
			avgFleetSize = float64(data.TotalFleetSize) / float64(data.Count)
		}
		fleetValueData = append(fleetValueData, FleetSizeAndValueData{
			Time:         t,
			AvgFleetSize: avgFleetSize,
			TotalValue:   data.TotalValue,
		})
	}

	// Sort the slice by time in ascending order
	sort.Slice(fleetValueData, func(i, j int) bool {
		return fleetValueData[i].Time.Before(fleetValueData[j].Time)
	})

	return fleetValueData
}

// getMondayOfISOWeek returns the date of the Monday for the given ISO year and week
func getMondayOfISOWeek(year int, week int) time.Time {
	// The first week of the year is the one that contains the first Thursday of the year
	// This aligns with ISO 8601
	// Start from January 4th, which is always in the first ISO week
	firstThursday := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)
	_, isoWeek := firstThursday.ISOWeek()

	// Calculate the difference in weeks
	weekDiff := week - isoWeek

	// Calculate the date of the Monday of the desired week
	desiredMonday := firstThursday.AddDate(0, 0, (weekDiff*7)+(-int(firstThursday.Weekday()-1)))

	// Adjust for time zones if necessary
	return time.Date(desiredMonday.Year(), desiredMonday.Month(), desiredMonday.Day(), 0, 0, 0, 0, time.UTC)
}
