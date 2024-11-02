package visuals

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

// Global variables
var (
	trackedCharacters []int
	orchestrator      *service.OrchestrateService
	logger            *logrus.Logger
)

// TemplateData holds all the data passed to the template
type TemplateData struct {
	// Chart data for MTD, YTD, Last Month
	MTDCharacterDamageData        template.JS
	YTDCharacterDamageData        template.JS
	LastMCharacterDamageData      template.JS
	MTDOurLossesValueData         template.JS
	YTDOurLossesValueData         template.JS
	LastMOurLossesValueData       template.JS
	MTDCharacterPerformanceData   template.JS
	YTDCharacterPerformanceData   template.JS
	LastMCharacterPerformanceData template.JS
	MTDOurShipsUsedData           template.JS
	YTDOurShipsUsedData           template.JS
	LastMOurShipsUsedData         template.JS
	MTDKillActivityData           template.JS
	YTDKillActivityData           template.JS
	LastMKillActivityData         template.JS
	MTDKillHeatmapData            template.JS
	YTDKillHeatmapData            template.JS
	LastMKillHeatmapData          template.JS
	MTDKillLossRatioData          template.JS
	YTDKillLossRatioData          template.JS
	LastMKillLossRatioData        template.JS
	MTDTopShipsKilledData         template.JS
	YTDTopShipsKilledData         template.JS
	LastMTopShipsKilledData       template.JS
	MTDVictimsByCorpData          template.JS
	YTDVictimsByCorpData          template.JS
	LastMVictimsByCorpData        template.JS
	MTDValueOverTimeData          template.JS
	YTDValueOverTimeData          template.JS
	LastMValueOverTimeData        template.JS
	MTDAverageFleetSizeData       template.JS
	YTDAverageFleetSizeData       template.JS
	LastMAverageFleetSizeData     template.JS
}

// Chart represents a single chart with its data preparation function and field prefix
type Chart struct {
	FieldPrefix string
	PrepareFunc func(*model.ChartData) interface{}
	Description string
}

// Define all charts with their corresponding preparation functions and field prefixes
var chartDefinitions = []Chart{
	{
		FieldPrefix: "CharacterDamageData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetDamageAndFinalBlows(cd)
		},
		Description: "Character Damage and Final Blows",
	},
	{
		FieldPrefix: "OurLossesValueData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetCombinedLossData(cd)
		},
		Description: "Our Losses",
	},
	{
		FieldPrefix: "CharacterPerformanceData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetCharacterPerformance(cd)
		},
		Description: "Character Performance",
	},
	{
		FieldPrefix: "OurShipsUsedData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetOurShipsUsed(cd)
		},
		Description: "Our Ships Used",
	},
	{
		FieldPrefix: "KillActivityData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetKillActivityOverTime(cd, "daily")
		},
		Description: "Kill Activity",
	},
	{
		FieldPrefix: "KillHeatmapData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetKillHeatmapData(cd)
		},
		Description: "Kill HeatMap",
	},
	{
		FieldPrefix: "KillLossRatioData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetKillLossRatioData(cd)
		},
		Description: "Kill Loss Ratio",
	},
	{
		FieldPrefix: "TopShipsKilledData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetTopShipsKilledData(cd)
		},
		Description: "Top Ships Killed",
	},
	{
		FieldPrefix: "VictimsByCorpData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetVictimsByCorp(cd)
		},
		Description: "Victims by Corp",
	},
	{
		FieldPrefix: "ValueOverTimeData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetValueOverTimeData(cd, "daily")
		},
		Description: "Value Over Time",
	},
	{
		FieldPrefix: "AverageFleetSizeData",
		PrepareFunc: func(cd *model.ChartData) interface{} {
			return GetAverageFleetSizeOverTime(cd, "daily")
		},
		Description: "Average Fleet Size Over Time",
	},
}

// TimeFrame represents the different time frames for the charts
type TimeFrame struct {
	Prefix string
	Data   *model.ChartData
}

// RenderSnippets prepares the template data and renders the template to a file
func RenderSnippets(orchestrateService *service.OrchestrateService, ytdChartData, lastMonthChartData, mtdChartData *model.ChartData, filePath string) error {
	orchestrator = orchestrateService
	logger = orchestrator.Logger

	// Fetch tracked characters from OrchestrateService
	ytdTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(ytdChartData.KillMails, &ytdChartData.ESIData)
	lastMTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(lastMonthChartData.KillMails, &lastMonthChartData.ESIData)
	mtdTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(mtdChartData.KillMails, &mtdChartData.ESIData)

	trackedCharacters = append(trackedCharacters, ytdTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, lastMTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, mtdTrackedCharacters...)

	orchestrator.Logger.Infof("there are %d tracked characters", len(trackedCharacters))

	data := TemplateData{}

	// Define time frames
	timeFrames := []TimeFrame{
		{"MTD", mtdChartData},
		{"YTD", ytdChartData},
		{"LastM", lastMonthChartData},
	}

	// Use reflection to dynamically set fields in TemplateData
	v := reflect.ValueOf(&data).Elem()

	for _, chart := range chartDefinitions {
		for _, tf := range timeFrames {
			// Construct the field name, e.g., "MTDCharacterDamageData"
			fieldName := fmt.Sprintf("%s%s", tf.Prefix, chart.FieldPrefix)
			preparedData, err := prepareData(tf.Data, chart.PrepareFunc, chart.Description)
			if err != nil {
				orchestrator.Logger.Errorf("Error preparing data for %s: %v", chart.Description, err)
				preparedData = template.JS("[]") // Fallback to empty array
			}

			// Set the field using reflection
			field := v.FieldByName(fieldName)
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(preparedData))
			} else {
				orchestrator.Logger.Errorf("Invalid field name: %s", fieldName)
				return fmt.Errorf("invalid field name: %s", fieldName)
			}
		}
	}

	// Render the template
	tmpl, err := template.New("tps.tmpl").ParseFiles(filepath.Join("static", "tmpl", "tps.tmpl"))
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
