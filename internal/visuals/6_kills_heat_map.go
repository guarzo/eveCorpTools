package visuals

import "github.com/guarzo/zkillanalytics/internal/model"

type HeatmapData struct {
	DayOfWeek int // 0 = Sunday, 6 = Saturday
	Hour      int // 0 - 23
	Kills     int
}

func GetKillHeatmapData(chartData *model.ChartData) [][]int {
	// Initialize a 7x24 matrix
	heatmap := make([][]int, 7)
	for i := range heatmap {
		heatmap[i] = make([]int, 24)
	}

	for _, km := range chartData.KillMails {
		timestamp := km.EsiKillMail.KillMailTime
		dayOfWeek := int(timestamp.Weekday())
		hour := timestamp.Hour()

		heatmap[dayOfWeek][hour]++
	}

	return heatmap
}
