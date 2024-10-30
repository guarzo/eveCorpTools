package visuals

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

var trackedCharacters []int
var orchestrator *service.OrchestrateService

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

func RenderSnippets(orchestrateService *service.OrchestrateService, ytdChartData, lastMonthChartData, mtdChartData *model.ChartData, filePath string) error {
	orchestrator = orchestrateService

	// Fetch tracked characters from OrchestrateService
	ytdTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(ytdChartData.KillMails, &ytdChartData.ESIData)
	lastMTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(lastMonthChartData.KillMails, &lastMonthChartData.ESIData)
	mtdTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(mtdChartData.KillMails, &mtdChartData.ESIData)

	trackedCharacters = append(trackedCharacters, ytdTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, lastMTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, mtdTrackedCharacters...)

	orchestrateService.Logger.Infof("there are %d tracked characters", len(trackedCharacters))

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
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

func prepareOurLossesValueData(chartData *model.ChartData) template.JS {
	data := GetOurLossesValue(chartData)
	if len(data) == 0 {
		fmt.Println("Warning: Our Losses Value Data is empty.")
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshalling OurLossesValueData: %v\n", err)
		return template.JS("[]") // Return empty array to prevent JS errors
	}
	return template.JS(jsonData)
}

func prepareCharacterPerformanceData(chartData *model.ChartData) template.JS {
	data := GetCharacterPerformance(chartData)
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

func prepareOurShipsUsedData(chartData *model.ChartData) template.JS {
	data := GetOurShipsUsed(chartData)
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

func prepareKillActivityData(chartData *model.ChartData) template.JS {
	data := GetKillActivityOverTime(chartData, "daily")
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

func prepareKillHeatmapData(chartData *model.ChartData) template.JS {
	data := GetKillHeatmapData(chartData)
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

func prepareKillLossRatioData(chartData *model.ChartData) template.JS {
	data := GetKillLossRatioData(chartData)
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

func prepareTopShipsKilledData(chartData *model.ChartData) template.JS {
	data := GetTopShipsKilled(chartData)
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

func prepareValueOverTimeData(chartData *model.ChartData) template.JS {
	data := GetValueOverTimeData(chartData, "daily")
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

func prepareKillCountData(chartData *model.ChartData) template.JS {
	data, _ := PrepareKillCountChartData(chartData)
	jsonData, _ := json.Marshal(data)
	return template.JS(jsonData)
}

// Implement the data preparation functions (similar to previous examples)
// Due to space constraints, I'll include just one as an example

// GetCharacterPerformance
func GetCharacterPerformance(chartData *model.ChartData) []CharacterKillData {
	characterStats := make(map[int]*CharacterKillData)

	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterID := attacker.CharacterID
			if characterID == 0 {
				continue
			}

			if !isOurCharacter(characterID) {
				continue
			}

			// Get character info
			characterInfo := chartData.CharacterInfos[characterID]
			//if characterInfo == nil {
			//	continue
			//}

			data, exists := characterStats[characterID]
			if !exists {
				data = &CharacterKillData{
					Name: characterInfo.Name,
				}
				characterStats[characterID] = data
			}
			data.KillCount++
			data.Points += km.ZKB.Points

			// Check for solo kill
			if len(km.EsiKillMail.Attackers) == 1 {
				data.SoloKills++
			}
		}
	}

	// Convert map to slice
	var result []CharacterKillData
	for _, data := range characterStats {
		result = append(result, *data)
	}

	// Sort by kill count
	sort.Slice(result, func(i, j int) bool {
		return result[i].KillCount > result[j].KillCount
	})

	return result
}

// Implement other data preparation functions similarly

// isOurCharacter helper function
func isOurCharacter(characterID int) bool {
	for _, id := range trackedCharacters {
		if id == characterID {
			return true
		}
	}
	return false
}
