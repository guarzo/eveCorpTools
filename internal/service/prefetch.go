// internal/service/prefetch.go

package service

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/persist"
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
	pf.saveCacheOnExit() // Register exit handler here instead of in Stop
	pf.Logger.Info("PrefetchService started.")
}

// run contains the main loop for prefetching.
func (pf *PrefetchService) run(ctx context.Context) {
	defer pf.wg.Done()
	pf.Logger.Info("PrefetchService is running.")

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	saveTicker := time.NewTicker(1 * time.Hour) // Save cache every hour
	defer saveTicker.Stop()

	// Initial prefetch immediately upon starting
	pf.prefetch(ctx)

	for {
		select {
		case <-ticker.C:
			pf.prefetch(ctx)
		case <-saveTicker.C:
			cacheFile := persist.GenerateCacheDataFileName()
			if err := pf.OrchestrateService.Cache.SaveToFile(cacheFile); err != nil {
				pf.Logger.Errorf("Failed to save cache periodically: %v", err)
			} else {
				pf.Logger.Infof("Cache saved periodically to %s", cacheFile)
			}
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
	pf.saveCacheOnExit()

	done := make(chan struct{})
	go func() {
		pf.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		pf.Logger.Info("PrefetchService stopped.")
	case <-time.After(30 * time.Second):
		pf.Logger.Error("PrefetchService did not stop within the timeout. Forcing exit.")
		os.Exit(1)
	}
}

// saveCacheOnExit sets up a signal handler to save the cache on program exit
func (pf *PrefetchService) saveCacheOnExit() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan // Wait for exit signal

		cache := pf.OrchestrateService.Cache
		cacheFile := persist.GenerateCacheDataFileName()
		if err := cache.SaveToFile(cacheFile); err != nil {
			log.Fatalf("Failed to save cache to file on exit: %v", err)
		} else {
			log.Printf("Cache successfully saved to %s on exit", cacheFile)
		}
		os.Exit(0) // Exit after saving
	}()
}
