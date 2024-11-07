package tps

import (
	"bytes"
	"context"
	"fmt"
	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/service"
	"github.com/guarzo/zkillanalytics/internal/visuals"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

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
	//case persist-trust.Config:
	//	configHandler(w)
	//	return nil
	default:
		orchestrator.Logger.Infof("Fetching data for %v", route)
		lastMonthData, err := fetchDataForSnippets(orchestrator, config.PreviousMonth)
		if err != nil {
			return err
		}
		mtdData, err := fetchDataForSnippets(orchestrator, config.MonthToDate)
		if err != nil {
			return err
		}
		orchestrator.Logger.Infof("Rendering chart for %v", route)
		return visuals.RenderCharts(orchestrator, chartData, lastMonthData, mtdData, filePath)
	}
}

func fetchDataForSnippets(orchestrator *service.OrchestrateService, dataMode config.DataMode) (*model.ChartData, error) {
	// Retrieve initial date range in string format
	startDateStr, endDateStr := persist.GetDateRange(dataMode)

	// Convert startDateStr and endDateStr to time.Time
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		orchestrator.Logger.Errorf("Invalid start date format: %v", err)
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		orchestrator.Logger.Errorf("Invalid end date format: %v", err)
		return nil, err
	}

	// Adjust date ranges explicitly for each DataMode
	if dataMode == config.PreviousMonth {
		// Set endDate to the last day of the previous month
		endDate = startDate.AddDate(0, 1, -1) // Moves to end of the start month
		endDateStr = endDate.Format("2006-01-02")
	} else if dataMode == config.MonthToDate {
		// Set endDate to today's date for MonthToDate
		endDate = time.Now()
		endDateStr = endDate.Format("2006-01-02")
	}

	// Log final adjusted dates for verification
	orchestrator.Logger.Infof("Fetching data for %v from %s to %s", dataMode, startDateStr, endDateStr)
	if dataMode == config.PreviousMonth {
		orchestrator.Logger.Infof("Expected Last Month Data: %s to %s", startDateStr, endDateStr)
	} else if dataMode == config.MonthToDate {
		orchestrator.Logger.Infof("Expected MTD Data: %s to %s", startDateStr, endDateStr)
	}

	// Fetch data with adjusted date range
	return orchestrator.GetAllData(context.TODO(), config.CorporationIDs, config.AllianceIDs, config.CharacterIDs, startDateStr, endDateStr)
}

func LoadingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loading page redirect")

	// Parse template
	tmplPath := filepath.Join("static", "tmpl", "loading.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		fmt.Printf("Error parsing template %s: %v\n", tmplPath, err)
		http.Error(w, "Loading Page Not Found", http.StatusNotFound)
		return
	}

	// Execute the template into a buffer first
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		fmt.Printf("Error executing template %s: %v\n", tmplPath, err)
		http.Error(w, "Failed to render loading page", http.StatusInternalServerError)
		return
	}

	// Set headers and write the buffer to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf.Bytes()); err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}

	// Explicitly flush the response if possible
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}
