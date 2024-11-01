package visuals

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

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
	// Existing data
	MTDKillCountData   template.JS
	YTDKillCountData   template.JS
	LastMKillCountData template.JS

	// Chart data for MTD, YTD, Last Month
	MTDCharacterDamageData   template.JS
	YTDCharacterDamageData   template.JS
	LastMCharacterDamageData template.JS

	MTDOurLossesValueData   template.JS
	YTDOurLossesValueData   template.JS
	LastMOurLossesValueData template.JS

	MTDCharacterPerformanceData   template.JS
	YTDCharacterPerformanceData   template.JS
	LastMCharacterPerformanceData template.JS

	MTDOurShipsUsedData   template.JS
	YTDOurShipsUsedData   template.JS
	LastMOurShipsUsedData template.JS

	MTDVictimsSunburstData   template.JS
	YTDVictimsSunburstData   template.JS
	LastMVictimsSunburstData template.JS

	MTDKillActivityData   template.JS
	YTDKillActivityData   template.JS
	LastMKillActivityData template.JS

	MTDKillHeatmapData   template.JS
	YTDKillHeatmapData   template.JS
	LastMKillHeatmapData template.JS

	MTDKillLossRatioData   template.JS
	YTDKillLossRatioData   template.JS
	LastMKillLossRatioData template.JS

	MTDTopShipsKilledData   template.JS
	YTDTopShipsKilledData   template.JS
	LastMTopShipsKilledData template.JS

	MTDValueOverTimeData   template.JS
	YTDValueOverTimeData   template.JS
	LastMValueOverTimeData template.JS
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

	// Prepare data for each chart and each time frame

	// 1. Damage Done and Final Blows
	data.MTDCharacterDamageData = prepareCharacterDamageData(mtdChartData)
	data.YTDCharacterDamageData = prepareCharacterDamageData(ytdChartData)
	data.LastMCharacterDamageData = prepareCharacterDamageData(lastMonthChartData)

	// 2. Our Losses Combined Chart
	data.MTDOurLossesValueData = prepareOurLossesValueData(mtdChartData)
	data.YTDOurLossesValueData = prepareOurLossesValueData(ytdChartData)
	data.LastMOurLossesValueData = prepareOurLossesValueData(lastMonthChartData)

	// 3. Character Performance Chart
	data.MTDCharacterPerformanceData = prepareCharacterPerformanceData(mtdChartData)
	data.YTDCharacterPerformanceData = prepareCharacterPerformanceData(ytdChartData)
	data.LastMCharacterPerformanceData = prepareCharacterPerformanceData(lastMonthChartData)

	// 4. Our Ships Used Chart
	data.MTDOurShipsUsedData = prepareOurShipsUsedData(mtdChartData)
	data.YTDOurShipsUsedData = prepareOurShipsUsedData(ytdChartData)
	data.LastMOurShipsUsedData = prepareOurShipsUsedData(lastMonthChartData)

	// 5. Victims Sunburst Chart
	data.MTDVictimsSunburstData = prepareVictimsSunburstData(mtdChartData)
	data.YTDVictimsSunburstData = prepareVictimsSunburstData(ytdChartData)
	data.LastMVictimsSunburstData = prepareVictimsSunburstData(lastMonthChartData)

	// 6. Kill Activity Over Time Chart
	data.MTDKillActivityData = prepareKillActivityData(mtdChartData)
	data.YTDKillActivityData = prepareKillActivityData(ytdChartData)
	data.LastMKillActivityData = prepareKillActivityData(lastMonthChartData)

	// 7. Kill Heatmap Chart
	data.MTDKillHeatmapData = prepareKillHeatmapData(mtdChartData)
	data.YTDKillHeatmapData = prepareKillHeatmapData(ytdChartData)
	data.LastMKillHeatmapData = prepareKillHeatmapData(lastMonthChartData)

	// 8. Kill-to-Loss Ratio Chart
	data.MTDKillLossRatioData = prepareKillLossRatioData(mtdChartData)
	data.YTDKillLossRatioData = prepareKillLossRatioData(ytdChartData)
	data.LastMKillLossRatioData = prepareKillLossRatioData(lastMonthChartData)

	// 9. Top Ships Killed Chart
	data.MTDTopShipsKilledData = prepareTopShipsKilledData(mtdChartData)
	data.YTDTopShipsKilledData = prepareTopShipsKilledData(ytdChartData)
	data.LastMTopShipsKilledData = prepareTopShipsKilledData(lastMonthChartData)

	// 10. Value Over Time Chart
	data.MTDValueOverTimeData = prepareValueOverTimeData(mtdChartData)
	data.YTDValueOverTimeData = prepareValueOverTimeData(ytdChartData)
	data.LastMValueOverTimeData = prepareValueOverTimeData(lastMonthChartData)

	// Prepare existing Kill Count Data
	data.MTDKillCountData = prepareKillCountData(mtdChartData)
	data.YTDKillCountData = prepareKillCountData(ytdChartData)
	data.LastMKillCountData = prepareKillCountData(lastMonthChartData)

	// Render the template
	tmpl, err := template.New("tps.tmpl").ParseFiles(filepath.Join("static", "tps.tmpl"))
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

// Helper functions to prepare data and marshal to JSON

func prepareCharacterDamageData(chartData *model.ChartData) template.JS {
	data := GetDamageAndFinalBlows(chartData)
	logger.Infof("Character Damage and Final Blows %v", data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling CharacterDamageData: %v", err)
		return template.JS("[]") // Return empty array to prevent JS errors
	}
	return template.JS(jsonData)
}

func prepareOurLossesValueData(chartData *model.ChartData) template.JS {
	data := GetOurLossesValue(chartData)
	logger.Infof("Our Losses %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling OurLossesValueData: %v", err)
		return template.JS("[]") // Return empty array to prevent JS errors
	}
	return template.JS(jsonData)
}

func prepareCharacterPerformanceData(chartData *model.ChartData) template.JS {
	data := GetCharacterPerformance(chartData)
	logger.Infof("Character Performance %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling CharacterPerformanceData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

func prepareOurShipsUsedData(chartData *model.ChartData) template.JS {
	data := GetOurShipsUsed(chartData)
	logger.Infof("Our Ships Used %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling OurShipsUsedData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

func prepareVictimsSunburstData(chartData *model.ChartData) template.JS {
	data := GetVictimsSunburst(chartData)
	logger.Infof("Victim Sunburst %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling VictimsSunburstData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

func prepareKillActivityData(chartData *model.ChartData) template.JS {
	data := GetKillActivityOverTime(chartData, "daily")
	logger.Infof("Kill Activity %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling KillActivityData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

func prepareKillHeatmapData(chartData *model.ChartData) template.JS {
	data := GetKillHeatmapData(chartData)
	logger.Infof("Kill HeatMap %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling KillHeatmapData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

func prepareKillLossRatioData(chartData *model.ChartData) template.JS {
	data := GetKillLossRatioData(chartData)
	logger.Infof("Kill Loss Ratio %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling KillLossRatioData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

func prepareTopShipsKilledData(chartData *model.ChartData) template.JS {
	data := GetTopShipsKilledData(chartData)
	logger.Infof("Top Ships KIlled %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling TopShipsKilledData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

func prepareValueOverTimeData(chartData *model.ChartData) template.JS {
	data := GetValueOverTimeData(chartData, "daily")
	logger.Infof("Value over time %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling ValueOverTimeData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

func prepareKillCountData(chartData *model.ChartData) template.JS {
	data, err := PrepareKillCountChartData(chartData)
	logger.Infof("Kill Count %v", data)

	if err != nil {
		orchestrator.Logger.Errorf("Error preparing KillCountData: %v", err)
		return template.JS("[]")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling KillCountData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}
