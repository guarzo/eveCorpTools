package visuals

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gambtho/zkillanalytics/internal/fetch"
	"github.com/gambtho/zkillanalytics/internal/model"
)

var trackedCharacters []int

func RenderSnippets(ytdChartData, lastMonthChartData, mtdChartData *model.ChartData, filePath string) error {
	ytdTrackedCharacters := fetch.GetTrackedCharacters(ytdChartData.KillMails, &ytdChartData.ESIData)
	lastMTrackedCharacters := fetch.GetTrackedCharacters(lastMonthChartData.KillMails, &lastMonthChartData.ESIData)
	mtdTrackedCharacters := fetch.GetTrackedCharacters(mtdChartData.KillMails, &mtdChartData.ESIData)

	trackedCharacters = append(trackedCharacters, ytdTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, lastMTrackedCharacters...)
	trackedCharacters = append(trackedCharacters, mtdTrackedCharacters...)

	data := struct {
		MTDChartHTML template.HTML
		YTDChartHTML template.HTML
		LMChartHTML  template.HTML
	}{
		MTDChartHTML: renderToHtml(RenderTopCharacters(mtdChartData), RenderOurShips(mtdChartData), RenderPointsPerCharacter(mtdChartData), RenderDamageDone(mtdChartData, "damage"), RenderDamageDone(mtdChartData, "blows"), RenderSolo(mtdChartData), RenderWeaponsByCharacter(mtdChartData), RenderVictims(mtdChartData), RenderTopShipsKilled(mtdChartData), RenderLostShipTypes(mtdChartData), RenderOurLossesValue(mtdChartData), RenderOurLossesCount(mtdChartData)),
		YTDChartHTML: renderToHtml(RenderTopCharacters(ytdChartData), RenderOurShips(ytdChartData), RenderPointsPerCharacter(ytdChartData), RenderDamageDone(ytdChartData, "damage"), RenderDamageDone(ytdChartData, "blows"), RenderSolo(ytdChartData), RenderWeaponsByCharacter(ytdChartData), RenderVictims(ytdChartData), RenderTopShipsKilled(ytdChartData), RenderLostShipTypes(ytdChartData), RenderOurLossesValue(ytdChartData), RenderOurLossesCount(ytdChartData)),
		LMChartHTML:  renderToHtml(RenderTopCharacters(lastMonthChartData), RenderOurShips(lastMonthChartData), RenderPointsPerCharacter(lastMonthChartData), RenderDamageDone(lastMonthChartData, "damage"), RenderDamageDone(lastMonthChartData, "blows"), RenderSolo(lastMonthChartData), RenderWeaponsByCharacter(lastMonthChartData), RenderVictims(lastMonthChartData), RenderTopShipsKilled(lastMonthChartData), RenderLostShipTypes(lastMonthChartData), RenderOurLossesValue(lastMonthChartData), RenderOurLossesCount(lastMonthChartData)),
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
