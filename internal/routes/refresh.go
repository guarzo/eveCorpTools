package routes

import (
	"context"
	"net/http"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/service"
)

func RefreshTPSHandler(orchestrateService *service.OrchestrateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := persist.DeleteFilesInDirectory(persist.GetChartsDirectory())
		if err != nil {
			orchestrateService.Logger.Infof("Using new charts directory: %v", err)
		}
		orchestrateService.Logger.Infof("Charts directory emptied for refresh")

		err = persist.DeleteCurrentMonthFile()
		if err != nil {
			orchestrateService.Logger.Errorf("Error removing current month data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		orchestrateService.Logger.Infof("removed current month data")

		begin, end := persist.GetDateRange(config.MonthToDate)
		_, err = orchestrateService.GetAllData(context.Background(), config.CorporationIDs, config.AllianceIDs, config.CharacterIDs, begin, end)
		if err != nil {
			orchestrateService.Logger.Errorf("Error fetching updated killmails: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orchestrateService.Logger.Infof("Updated current month data")

		// Set the Content-Type header to indicate plain text response
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// Write the success status code
		w.WriteHeader(http.StatusOK)

		// Write the success message to the response body
		_, writeErr := w.Write([]byte("Refresh successful"))
		if writeErr != nil {
			// Log the error if writing the response fails
			orchestrateService.Logger.Errorf("Failed to write response: %v", writeErr)
		}
	}
}
