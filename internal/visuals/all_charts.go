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

	MTDVictimsByCorpData   template.JS
	YTDVictimsByCorpData   template.JS
	LastMVictimsByCorpData template.JS

	MTDValueOverTimeData   template.JS
	YTDValueOverTimeData   template.JS
	LastMValueOverTimeData template.JS

	MTDAverageFleetSizeData   template.JS
	YTDAverageFleetSizeData   template.JS
	LastMAverageFleetSizeData template.JS
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
	data.MTDOurLossesValueData = prepareCombinedLossData(mtdChartData)
	data.YTDOurLossesValueData = prepareCombinedLossData(ytdChartData)
	data.LastMOurLossesValueData = prepareCombinedLossData(lastMonthChartData)

	// 3. Character Performance Chart
	data.MTDCharacterPerformanceData = prepareCharacterPerformanceData(mtdChartData)
	data.YTDCharacterPerformanceData = prepareCharacterPerformanceData(ytdChartData)
	data.LastMCharacterPerformanceData = prepareCharacterPerformanceData(lastMonthChartData)

	// 4. Our Ships Used Chart
	data.MTDOurShipsUsedData = prepareOurShipsUsedData(mtdChartData)
	data.YTDOurShipsUsedData = prepareOurShipsUsedData(ytdChartData)
	data.LastMOurShipsUsedData = prepareOurShipsUsedData(lastMonthChartData)

	// 5. Kill Activity Over Time Chart
	data.MTDKillActivityData = prepareKillActivityData(mtdChartData)
	data.YTDKillActivityData = prepareKillActivityData(ytdChartData)
	data.LastMKillActivityData = prepareKillActivityData(lastMonthChartData)

	// 6. Kill Heatmap Chart
	data.MTDKillHeatmapData = prepareKillHeatmapData(mtdChartData)
	data.YTDKillHeatmapData = prepareKillHeatmapData(ytdChartData)
	data.LastMKillHeatmapData = prepareKillHeatmapData(lastMonthChartData)

	// 7. Kill-to-Loss Ratio Chart
	data.MTDKillLossRatioData = prepareKillLossRatioData(mtdChartData)
	data.YTDKillLossRatioData = prepareKillLossRatioData(ytdChartData)
	data.LastMKillLossRatioData = prepareKillLossRatioData(lastMonthChartData)

	// 8. Top Ships Killed Chart
	data.MTDTopShipsKilledData = prepareTopShipsKilledData(mtdChartData)
	data.YTDTopShipsKilledData = prepareTopShipsKilledData(ytdChartData)
	data.LastMTopShipsKilledData = prepareTopShipsKilledData(lastMonthChartData)

	// 9. Victims by Corp Chart
	data.MTDVictimsByCorpData = prepareVictimsByCorp(mtdChartData)
	data.YTDVictimsByCorpData = prepareVictimsByCorp(ytdChartData)
	data.LastMVictimsByCorpData = prepareVictimsByCorp(lastMonthChartData)

	// 10. Value Over Time Chart
	data.MTDAverageFleetSizeData = prepareAverageFleetSizeData(mtdChartData)
	data.YTDAverageFleetSizeData = prepareAverageFleetSizeData(ytdChartData)
	data.LastMAverageFleetSizeData = prepareAverageFleetSizeData(lastMonthChartData)

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

func prepareCombinedLossData(chartData *model.ChartData) template.JS {
	data := GetCombinedLossData(chartData)
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

func prepareVictimsByCorp(chartData *model.ChartData) template.JS {
	data := GetVictimsByCorp(chartData)
	logger.Infof("Victims by corp %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling VictimsByCorpData: %v", err)
		return "[]"
	}
	return template.JS(jsonData)
}

func prepareAverageFleetSizeData(chartData *model.ChartData) template.JS {
	data := GetAverageFleetSizeOverTime(chartData, "daily")
	logger.Infof("Average FleetSize over time %v", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		orchestrator.Logger.Errorf("Error marshalling AverageFleetSizeOverTimeData: %v", err)
		return template.JS("[]")
	}
	return template.JS(jsonData)
}

// CharacterKillData holds the data for character kill counts
type CharacterKillData struct {
	CharacterID int
	KillCount   int
	Name        string
	Points      int
	SoloKills   int
}

type ChartJSData struct {
	Labels   []string         `json:"labels"`
	Datasets []ChartJSDataset `json:"datasets"`
}

type ChartJSDataset struct {
	Label           string   `json:"label"`
	Data            []int    `json:"data"`
	BackgroundColor []string `json:"backgroundColor"`
}
