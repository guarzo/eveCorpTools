package visuals

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

// In your visuals package
type ChartJSData struct {
	Labels   []string         `json:"labels"`
	Datasets []ChartJSDataset `json:"datasets"`
}

type ChartJSDataset struct {
	Label           string   `json:"label"`
	Data            []int    `json:"data"`
	BackgroundColor []string `json:"backgroundColor"`
}

var trackedCharacters []int

// RenderSnippets generates and renders chart snippets based on provided chart data.
func RenderSnippets(orchestrateService *service.OrchestrateService, ytdChartData, lastMonthChartData, mtdChartData *model.ChartData, filePath string) error {
	// Fetch tracked characters from OrchestrateService
	ytdTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(ytdChartData.KillMails, &ytdChartData.ESIData)
	lastMTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(lastMonthChartData.KillMails, &lastMonthChartData.ESIData)
	mtdTrackedCharacters := orchestrateService.GetTrackedCharactersFromKillMails(mtdChartData.KillMails, &mtdChartData.ESIData)

	trackedCharacters = append(trackedCharacters, ytdTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, lastMTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, mtdTrackedCharacters...)

	orchestrateService.Logger.Infof("there are %d tracked characters", len(trackedCharacters))

	mtdKillCountData, err := PrepareKillCountChartData(mtdChartData)
	if err != nil {
		return fmt.Errorf("failed to prepare MTD Kill Count data: %w", err)
	}
	mtdKillCountDataJSON, err := json.Marshal(mtdKillCountData)
	if err != nil {
		return fmt.Errorf("failed to marshal MTD Kill Count data: %w", err)
	}

	ytdKillCountData, err := PrepareKillCountChartData(ytdChartData)
	if err != nil {
		return fmt.Errorf("failed to prepare ytd Kill Count data: %w", err)
	}
	ytdKillCountDataJSON, err := json.Marshal(ytdKillCountData)
	if err != nil {
		return fmt.Errorf("failed to marshal ytd Kill Count data: %w", err)
	}

	lastMKillCountData, err := PrepareKillCountChartData(lastMonthChartData)
	if err != nil {
		return fmt.Errorf("failed to prepare lastM Kill Count data: %w", err)
	}
	lastMKillCountDataJSON, err := json.Marshal(lastMKillCountData)
	if err != nil {
		return fmt.Errorf("failed to marshal lastM Kill Count data: %w", err)
	}

	data := struct {
		MTDKillCountData   template.JS
		YTDKillCountData   template.JS
		LastMKillCountData template.JS
	}{
		MTDKillCountData:   template.JS(mtdKillCountDataJSON),
		YTDKillCountData:   template.JS(ytdKillCountDataJSON),
		LastMKillCountData: template.JS(lastMKillCountDataJSON),
	}

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

type ChartData struct {
	ChartHTML template.HTML
}
