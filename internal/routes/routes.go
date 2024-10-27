package routes

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"

	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/service"
)

// ServeRoute is an HTTP handler that generates a bar chart based on the mode
func ServeRoute(route config.Route, orchestrateService *service.OrchestrateService) http.HandlerFunc {
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
			fmt.Println(fmt.Sprintf("Serving existing chart for %s from %s to %s based on mode %s", config.RouteToString[route], startDate, endDate, modeStr))
			http.ServeFile(w, r, filePath)
			return
		}

		orchestrateService.Logger.Infof("Creating chart for %s from %s to %s based on mode %s", config.RouteToString[route], startDate, endDate, modeStr)

		corporations := orchestrateService.GetTrackedCorporations()
		alliances := orchestrateService.GetTrackedAlliances()
		characters := orchestrateService.GetTrackedCharacters()

		// Fetch data and create chart using OrchestrateService
		chartData, err := orchestrateService.GetAllData(r.Context(), corporations, alliances, characters, startDate, endDate)
		if err != nil {
			if err.Error() == "another GetAllData operation is in progress" {
				orchestrateService.Logger.Warn("Another GetAllData operation is in progress")
				LoadingHandler(w, r)
			} else {
				orchestrateService.Logger.Errorf("Error fetching detailed killmails: %v", err)
				http.Error(w, fmt.Sprintf("Error fetching detailed killmails: %s", err), http.StatusInternalServerError)
			}
			return
		}

		if err := generateChart(orchestrateService, route, chartData, filePath, w); err != nil {
			http.Error(w, fmt.Sprintf("Error creating bar chart: %s", err), http.StatusInternalServerError)
			return
		}

		// Serve the HTML file
		http.ServeFile(w, r, filePath)
	}
}
