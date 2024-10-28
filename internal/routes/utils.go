package routes

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/service"
	"github.com/guarzo/zkillanalytics/internal/visuals"
)

//func CreateCorporationMap(client *http.Client, ids []int) (map[int]model.Namer, error) {
//	corporationMap := make(map[int]model.Namer)
//
//	for _, id := range ids {
//		info, err := esi.GetCorporationInfo(client, id)
//		if err != nil {
//			return nil, err
//		}
//		corporationMap[id] = info
//	}
//
//	return corporationMap, nil
//}
//
//func CreateAllianceMap(client *http.Client, ids []int) (map[int]model.Namer, error) {
//	allianceMap := make(map[int]model.Namer)
//
//	for _, id := range ids {
//		info, err := esi.GetAllianceInfo(client, id)
//		if err != nil {
//			return nil, err
//		}
//		allianceMap[id] = info
//	}
//
//	return allianceMap, nil
//}
//
//func CreateCharacterMap(client *http.Client, ids []int) (map[int]model.Namer, error) {
//	characterMap := make(map[int]model.Namer)
//
//	for _, id := range ids {
//		info, err := esi.GetCharacterInfo(client, id)
//		if err != nil {
//			return nil, err
//		}
//		characterMap[id] = info
//	}
//
//	return characterMap, nil
//}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("static", "404.tmpl"))
	if err != nil {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, nil)
}

func LoadingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loading page redirect")

	// Parse template
	tmplPath := filepath.Join("static", "loading.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		fmt.Printf("Error parsing template %s: %v\n", tmplPath, err)
		http.Error(w, "Loading Page Not Found", http.StatusNotFound)
		return
	}

	// Set headers and write response
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Error executing template %s: %v\n", tmplPath, err)
		http.Error(w, "Failed to render loading page", http.StatusInternalServerError)
		return
	}

	// Explicitly flush the response
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func getDataMode(modeStr, lastPart string) config.DataMode {
	dataMode, ok := config.StringToDataMode[modeStr]
	if !ok {
		dataMode = config.YearToDate
	}

	if lastPart == "lastMonth" {
		dataMode = config.PreviousMonth
	}
	if lastPart == "currentMonth" {
		dataMode = config.MonthToDate
	}

	return dataMode
}

func generateFilePath(dir string, route config.Route, startDate, endDate string) string {
	return persist.GenerateChartFileName(dir, config.RouteToString[route], startDate, endDate,
		persist.HashParams(persist.IntSliceToString(config.CorporationIDs)+persist.IntSliceToString(config.AllianceIDs)+persist.IntSliceToString(config.CharacterIDs)))
}

func generateChart(orchestrator *service.OrchestrateService, route config.Route, chartData *model.ChartData, filePath string, w http.ResponseWriter) error {
	fmt.Println("Generating chart for", config.RouteToString[route])
	switch route {
	//case persist.Config:
	//	configHandler(w)
	//	return nil
	default:
		lastMonthData, err := fetchDataForSnippets(orchestrator, config.PreviousMonth)
		if err != nil {
			return err
		}
		mtdData, err := fetchDataForSnippets(orchestrator, config.MonthToDate)
		if err != nil {
			return err
		}
		return visuals.RenderSnippets(orchestrator, chartData, lastMonthData, mtdData, filePath)
	}
}

func fetchDataForSnippets(orchestrator *service.OrchestrateService, dataMode config.DataMode) (*model.ChartData, error) {
	startDate, endDate := persist.GetDateRange(dataMode)
	return orchestrator.GetAllData(context.TODO(), config.CorporationIDs, config.AllianceIDs, config.CharacterIDs, startDate, endDate)
}
