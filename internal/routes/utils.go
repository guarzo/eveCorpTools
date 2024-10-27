package routes

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/service"
	"github.com/gambtho/zkillanalytics/internal/visuals"
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
	fmt.Println("loading page redirect`")
	tmpl, err := template.ParseFiles(filepath.Join("static", "loading.tmpl"))
	if err != nil {
		http.Error(w, "Loading Page Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func getDataMode(modeStr, lastPart string) persist.DataMode {
	dataMode, ok := persist.StringToDataMode[modeStr]
	if !ok {
		dataMode = persist.YearToDate
	}

	if lastPart == "lastMonth" {
		dataMode = persist.PreviousMonth
	}
	if lastPart == "currentMonth" {
		dataMode = persist.MonthToDate
	}

	return dataMode
}

func generateFilePath(dir string, route persist.Route, startDate, endDate string) string {
	return persist.GenerateChartFileName(dir, persist.RouteToString[route], startDate, endDate,
		persist.HashParams(persist.IntSliceToString(persist.CorporationIDs)+persist.IntSliceToString(persist.AllianceIDs)+persist.IntSliceToString(persist.CharacterIDs)))
}

func generateChart(orchestrator *service.OrchestrateService, route persist.Route, chartData *model.ChartData, filePath string, w http.ResponseWriter) error {
	fmt.Println("Generating chart for", persist.RouteToString[route])
	switch route {
	//case persist.Config:
	//	configHandler(w)
	//	return nil
	default:
		lastMonthData, err := fetchDataForSnippets(orchestrator, persist.PreviousMonth)
		if err != nil {
			return err
		}
		mtdData, err := fetchDataForSnippets(orchestrator, persist.MonthToDate)
		if err != nil {
			return err
		}
		return visuals.RenderSnippets(orchestrator, chartData, lastMonthData, mtdData, filePath)
	}
}

func fetchDataForSnippets(orchestrator *service.OrchestrateService, dataMode persist.DataMode) (*model.ChartData, error) {
	startDate, endDate := persist.GetDateRange(dataMode)
	return orchestrator.GetAllData(context.TODO(), persist.CorporationIDs, persist.AllianceIDs, persist.CharacterIDs, startDate, endDate)
}
