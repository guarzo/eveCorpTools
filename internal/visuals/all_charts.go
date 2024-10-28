package visuals

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

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

	data := struct {
		MTDChartHTML template.HTML
		YTDChartHTML template.HTML
		LMChartHTML  template.HTML
	}{
		MTDChartHTML: renderToHtml(RenderTopCharacters(mtdChartData), RenderOurShips(orchestrateService, mtdChartData), RenderPointsPerCharacter(mtdChartData), RenderDamageDone(mtdChartData, "damage"), RenderDamageDone(mtdChartData, "blows"), RenderSolo(mtdChartData), RenderWeaponsByCharacter(orchestrateService, mtdChartData), RenderVictims(mtdChartData), RenderTopShipsKilled(orchestrateService, mtdChartData), RenderLostShipTypes(orchestrateService, mtdChartData), RenderOurLossesValue(orchestrateService, mtdChartData), RenderOurLossesCount(orchestrateService, mtdChartData)),
		YTDChartHTML: renderToHtml(RenderTopCharacters(ytdChartData), RenderOurShips(orchestrateService, ytdChartData), RenderPointsPerCharacter(ytdChartData), RenderDamageDone(ytdChartData, "damage"), RenderDamageDone(ytdChartData, "blows"), RenderSolo(ytdChartData), RenderWeaponsByCharacter(orchestrateService, ytdChartData), RenderVictims(ytdChartData), RenderTopShipsKilled(orchestrateService, ytdChartData), RenderLostShipTypes(orchestrateService, ytdChartData), RenderOurLossesValue(orchestrateService, ytdChartData), RenderOurLossesCount(orchestrateService, ytdChartData)),
		LMChartHTML:  renderToHtml(RenderTopCharacters(lastMonthChartData), RenderOurShips(orchestrateService, lastMonthChartData), RenderPointsPerCharacter(lastMonthChartData), RenderDamageDone(lastMonthChartData, "damage"), RenderDamageDone(lastMonthChartData, "blows"), RenderSolo(lastMonthChartData), RenderWeaponsByCharacter(orchestrateService, lastMonthChartData), RenderVictims(lastMonthChartData), RenderTopShipsKilled(orchestrateService, lastMonthChartData), RenderLostShipTypes(orchestrateService, lastMonthChartData), RenderOurLossesValue(orchestrateService, lastMonthChartData), RenderOurLossesCount(orchestrateService, lastMonthChartData)),
	}

	tmpl, err := template.New("chart.tmpl").ParseFiles(filepath.Join("static", "chart.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Save the chart to an HTML file
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
