package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gambtho/zkillanalytics/internal/api/esi"
	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/service"
	"github.com/gambtho/zkillanalytics/internal/visuals"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func WaitForMutexAndCallFunction(retry int) bool {
	count := 0
	for count < retry {
		if service.FetchAllMutex.TryLock() {
			service.FetchAllMutex.Unlock()
			return true
		}
		count++
		time.Sleep(1 * time.Second) // Wait for 1 second before retrying
	}
	return false
}

func CreateCorporationMap(client *http.Client, ids []int) (map[int]model.Namer, error) {
	corporationMap := make(map[int]model.Namer)

	for _, id := range ids {
		info, err := esi.GetCorporationInfo(client, id)
		if err != nil {
			return nil, err
		}
		corporationMap[id] = info
	}

	return corporationMap, nil
}

func CreateAllianceMap(client *http.Client, ids []int) (map[int]model.Namer, error) {
	allianceMap := make(map[int]model.Namer)

	for _, id := range ids {
		info, err := esi.GetAllianceInfo(client, id)
		if err != nil {
			return nil, err
		}
		allianceMap[id] = info
	}

	return allianceMap, nil
}

func CreateCharacterMap(client *http.Client, ids []int) (map[int]model.Namer, error) {
	characterMap := make(map[int]model.Namer)

	for _, id := range ids {
		info, err := esi.GetCharacterInfo(client, id)
		if err != nil {
			return nil, err
		}
		characterMap[id] = info
	}

	return characterMap, nil
}

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

func getHttpClient() *http.Client {
	return &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("User-Agent", "guarzo.eve@gmail.com")
			req.Header.Set("Accept-Encoding", "gzip")
			return http.DefaultTransport.RoundTrip(req)
		}),
	}
}

func generateChart(route persist.Route, chartData *model.ChartData, filePath string, client *http.Client, w http.ResponseWriter) error {
	fmt.Println("Generating chart for", persist.RouteToString[route])
	switch route {
	case persist.Config:
		configHandler(w)
		return nil
	default:
		lastMonthData, err := fetchDataForSnippets(client, persist.PreviousMonth)
		if err != nil {
			return err
		}
		mtdData, err := fetchDataForSnippets(client, persist.MonthToDate)
		if err != nil {
			return err
		}
		return visuals.RenderSnippets(chartData, lastMonthData, mtdData, filePath)
	}
}

func fetchDataForSnippets(client *http.Client, dataMode persist.DataMode) (*model.ChartData, error) {
	startDate, endDate := persist.GetDateRange(dataMode)
	return service.GetAllData(client, persist.CorporationIDs, persist.AllianceIDs, persist.CharacterIDs, startDate, endDate)
}
