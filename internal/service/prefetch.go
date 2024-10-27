// internal/service/prefetch.go

package service

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

// PrefetchService handles scheduled data prefetching.
type PrefetchService struct {
	OrchestrateService *OrchestrateService

	// WaitGroup to track ongoing prefetch operations
	wg sync.WaitGroup

	Logger *logrus.Logger
}

// NewPrefetchService initializes and returns a new PrefetchService instance.
func NewPrefetchService(os *OrchestrateService, logger *logrus.Logger) *PrefetchService {
	return &PrefetchService{
		OrchestrateService: os,
		Logger:             logger,
	}
}

// Start begins the prefetching process.
func (pf *PrefetchService) Start(ctx context.Context) {
	pf.wg.Add(1)
	go pf.run(ctx)
	pf.Logger.Info("PrefetchService started.")
}

// run contains the main loop for prefetching.
func (pf *PrefetchService) run(ctx context.Context) {
	defer pf.wg.Done()
	pf.Logger.Info("PrefetchService is running.")

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Initial prefetch immediately upon starting
	pf.prefetch(ctx)

	for {
		select {
		case <-ticker.C:
			pf.prefetch(ctx)
		case <-ctx.Done():
			pf.Logger.Info("PrefetchService received context cancellation.")
			return
		}
	}
}

func (pf *PrefetchService) prefetch(ctx context.Context) {
	pf.Logger.Info("Starting prefetch operation.")
	begin, end := persist.GetDateRange(config.YearToDate)
	pf.Logger.Infof("Prefetching data for %s to %s...", begin, end)

	// Use a derived context with a timeout to prevent indefinite blocking
	prefetchCtx, cancel := context.WithTimeout(ctx, 23*time.Hour)
	defer cancel()

	pf.Logger.Info("Calling GetAllData...")
	chartData, err := pf.OrchestrateService.GetAllData(prefetchCtx, config.CorporationIDs, config.AllianceIDs, config.CharacterIDs, begin, end)
	if err != nil {
		pf.Logger.Errorf("Error fetching detailed killmails: %v", err)
		return
	}
	pf.Logger.Infof("Prefetch completed successfully with %d killmails.", len(chartData.KillMails))
}

// Stop gracefully waits for the prefetching process to complete with a timeout.
func (pf *PrefetchService) Stop() {
	pf.Logger.Info("Waiting for PrefetchService to stop...")

	done := make(chan struct{})
	go func() {
		pf.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		pf.Logger.Info("PrefetchService stopped.")
	case <-time.After(1 * time.Second):
		pf.Logger.Error("PrefetchService did not stop within the timeout. Forcing exit.")
		os.Exit(1)
	}
}
