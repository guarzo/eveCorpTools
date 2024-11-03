// internal/visuals/render_charts.go

package visuals

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

// Global variables
var (
	orchestrator *service.OrchestrateService
	logger       *logrus.Logger
)

// TemplateData holds all the data passed to the template
type TemplateData struct {
	TimeFrames []TimeFrameData
}

// TimeFrameData represents data for a specific time frame (MTD, YTD, LastM)
type TimeFrameData struct {
	Name   string       // e.g., "MTD", "YTD", "LastM"
	Charts []ChartEntry // Slice of charts for this time frame
}

// ChartEntry represents a single chart's data
type ChartEntry struct {
	Name string      // e.g., "Character Damage and Final Blows"
	ID   string      // e.g., "characterDamageAndFinalBlowsChart_MTD"
	Data template.JS // JSON data for the chart
	Type string      // e.g., "bar", "line", "matrix", "wordCloud"
}

// Chart represents a single chart with its data preparation function
type Chart struct {
	FieldPrefix string
	PrepareFunc func(*model.ChartData) interface{}
	Description string
	Type        string // e.g., "bar", "line", "matrix", "wordCloud"
}

// Define all charts with their corresponding preparation functions and field prefixes
var chartDefinitions = []Chart{
	{
		FieldPrefix: "CharacterDamageData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetDamageAndFinalBlows(cd)
		},
		Description: "Character Damage and Final Blows",
		Type:        "bar",
	},
	{
		FieldPrefix: "CharacterPerformanceData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetCharacterPerformance(cd)
		},
		Description: "Character Performance",
		Type:        "bar",
	},
	{
		FieldPrefix: "OurShipsUsedData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetOurShipsUsed(cd)
		},
		Description: "Our Ships Used",
		Type:        "bar",
	},
	{
		FieldPrefix: "KillActivityData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetKillActivityOverTime(cd, "daily")
		},
		Description: "Kill Activity Over Time",
		Type:        "line",
	},
	{
		FieldPrefix: "KillHeatmapData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetKillHeatmapData(cd)
		},
		Description: "Kills Heatmap",
		Type:        "matrix",
	},
	{
		FieldPrefix: "KillLossRatioData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetKillLossRatioData(cd)
		},
		Description: "Kill-to-Loss Ratio",
		Type:        "bar",
	},
	{
		FieldPrefix: "TopShipsKilledData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetTopShipsKilledData(cd)
		},
		Description: "Top Ships Killed",
		Type:        "wordCloud",
	},
	{
		FieldPrefix: "VictimsByCorpData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetVictimsByCorp(cd)
		},
		Description: "Victims by Corporation",
		Type:        "bar",
	},
	{
		FieldPrefix: "FleetSizeAndValueData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetFleetSizeAndValueData(cd, "daily")
		},
		Description: "Fleet Size and Value Killed Over Time",
		Type:        "line",
	},
	{
		FieldPrefix: "CombinedLossesData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetCombinedLossData(cd)
		},
		Description: "Combined Losses",
		Type:        "bar",
	},
}

// RenderCharts prepares the template data and renders the template to a file
func RenderCharts(orchestrateService *service.OrchestrateService, ytdChartData, lastMonthChartData, mtdChartData *model.ChartData, filePath string) error {
	orchestrator = orchestrateService
	logger = orchestrator.Logger

	// Fetch tracked characters from OrchestrateService
	ytdTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(ytdChartData.KillMails, &ytdChartData.ESIData)
	lastMTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(lastMonthChartData.KillMails, &lastMonthChartData.ESIData)
	mtdTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(mtdChartData.KillMails, &mtdChartData.ESIData)

	trackedCharacters := append(ytdTrackedCharacters, lastMTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, mtdTrackedCharacters...)

	orchestrator.Logger.Infof("there are %d tracked characters", len(trackedCharacters))

	data := TemplateData{
		TimeFrames: []TimeFrameData{
			{
				Name:   "MTD",
				Charts: []ChartEntry{},
			},
			{
				Name:   "LastM",
				Charts: []ChartEntry{},
			},
			{
				Name:   "YTD",
				Charts: []ChartEntry{},
			},
		},
	}

	// Define time frames and associate chart data
	timeFrames := []struct {
		Name string
		Data *model.ChartData
	}{
		{"MTD", mtdChartData},
		{"LastM", lastMonthChartData},
		{"YTD", ytdChartData},
	}

	// Populate TemplateData
	for _, tf := range timeFrames {
		for _, chart := range chartDefinitions {
			// Prepare data
			preparedData, err := prepareData(tf.Data, chart.PrepareFunc, chart.Description)
			if err != nil {
				orchestrator.Logger.Errorf("Error preparing data for %s: %v", chart.Description, err)
				preparedData = template.JS("[]") // Fallback to empty array
			}

			// Generate unique canvas ID based on Description and Timeframe
			chartID := fmt.Sprintf("%sChart_%s", toLowerCamelCase(chart.Description), tf.Name)

			// Append to charts for this time frame
			for i, finalTF := range data.TimeFrames {
				if finalTF.Name == tf.Name {
					data.TimeFrames[i].Charts = append(data.TimeFrames[i].Charts, ChartEntry{
						Name: chart.Description,
						ID:   chartID,
						Data: preparedData,
						Type: chart.Type,
					})
					break
				}
			}
		}
	}

	// Create a template.FuncMap with the toLower function
	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
	}

	// Render the template with the FuncMap
	tmpl, err := template.New("tps.tmpl").Funcs(funcMap).ParseFiles(filepath.Join("static", "tmpl", "tps.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}

// Helper function to convert description to lowerCamelCase for IDs
func toLowerCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	// Example: "Character Damage and Final Blows" -> "characterDamageAndFinalBlows"
	result := ""
	capitalizeNext := false
	for i, r := range s {
		if i == 0 {
			result += string(unicode.ToLower(r))
		} else if r == ' ' || r == '-' || r == '_' {
			capitalizeNext = true
		} else {
			if capitalizeNext {
				result += string(unicode.ToUpper(r))
				capitalizeNext = false
			} else {
				result += string(r)
			}
		}
	}
	return result
}

// Generic helper function to prepare data
func prepareData(chartData *model.ChartData, getDataFunc func(*model.ChartData) interface{}, description string) (template.JS, error) {
	data := getDataFunc(chartData)
	logger.Infof("%s: %v", description, data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling %s: %v", description, err)
		return "[]", err
	}
	return template.JS(jsonData), nil
}
