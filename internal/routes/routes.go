package routes

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"

	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/service"
)

// ServeRoute is an HTTP handler that generates a bar chart based on the mode
func ServeRoute(route persist.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse URL parameters to determine data mode
		vars := mux.Vars(r)
		modeStr := vars["mode"]
		lastPart := path.Base(r.URL.Path)

		dataMode := getDataMode(modeStr, lastPart)

		// Determine start and end dates
		startDate, endDate := persist.GetDateRange(dataMode)
		dir := persist.GetChartsDirectory()

		// Ensure the charts directory exists in the project root
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create charts directory: %s", err), http.StatusInternalServerError)
			return
		}

		filePath := generateFilePath(dir, route, startDate, endDate)

		// Check if the file already exists
		if _, err := os.Stat(filePath); err == nil {
			fmt.Println(fmt.Sprintf("Serving existing chart for %s from %s to %s based on mode %s", persist.RouteToString[route], startDate, endDate, modeStr))
			http.ServeFile(w, r, filePath)
			return
		}

		fmt.Println(fmt.Sprintf("Creating chart for %s from %s to %s based on mode %s", persist.RouteToString[route], startDate, endDate, modeStr))
		// Fetch data and create chart
		client := getHttpClient()

		if !WaitForMutexAndCallFunction(5) {
			fmt.Println("Failed to acquire mutex")
			LoadingHandler(w, r)
		}
		service.FetchAllMutex.Lock()
		defer service.FetchAllMutex.Unlock()

		chartData, err := service.GetAllData(client, persist.CorporationIDs, persist.AllianceIDs, persist.CharacterIDs, startDate, endDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching detailed killmails: %s", err), http.StatusInternalServerError)
			return
		}

		if err := generateChart(route, chartData, filePath, client, w); err != nil {
			http.Error(w, fmt.Sprintf("Error creating bar chart: %s", err), http.StatusInternalServerError)
			return
		}

		// Serve the HTML file
		http.ServeFile(w, r, filePath)
	}
}
