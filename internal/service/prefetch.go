package service

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gambtho/zkillanalytics/internal/persist"
)

var FetchAllMutex sync.Mutex

// PrefetchService handles scheduled data prefetching.
type PrefetchService struct {
	OrchestrateService *OrchestrateService

	// Control channels
	stopChan  chan struct{}
	stopped   chan struct{}
	runningMu sync.Mutex
	running   bool

	// WaitGroup to track ongoing prefetch operations
	wg sync.WaitGroup

	Logger *logrus.Logger
}

// NewPrefetchService initializes and returns a new PrefetchService instance.
func NewPrefetchService(os *OrchestrateService, logger *logrus.Logger) *PrefetchService {
	return &PrefetchService{
		OrchestrateService: os,
		stopChan:           make(chan struct{}),
		stopped:            make(chan struct{}),
		Logger:             logger,
	}
}

// Start begins the prefetching process.
// It ensures that only one instance runs at a time.
func (pf *PrefetchService) Start(ctx context.Context) {
	pf.runningMu.Lock()
	defer pf.runningMu.Unlock()

	if pf.running {
		pf.Logger.Info("PrefetchService is already running. Start request ignored.")
		return
	}

	pf.running = true
	pf.wg.Add(1)
	go pf.run(ctx)
	pf.Logger.Info("PrefetchService started.")
}

// run contains the main loop for prefetching.
func (pf *PrefetchService) run(ctx context.Context) {
	defer func() {
		pf.runningMu.Lock()
		pf.running = false
		pf.runningMu.Unlock()
		pf.wg.Done()
		close(pf.stopped)
		pf.Logger.Info("PrefetchService stopped.")
	}()

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Initial prefetch immediately upon starting
	pf.prefetch(ctx)

	for {
		select {
		case <-ticker.C:
			pf.prefetch(ctx)
		case <-pf.stopChan:
			pf.Logger.Info("PrefetchService received stop signal.")
			return
		case <-ctx.Done():
			pf.Logger.Info("PrefetchService received context cancellation.")
			return
		}
	}
}

// prefetch executes the data fetching logic.
func (pf *PrefetchService) prefetch(ctx context.Context) {
	begin, end := persist.GetDateRange(persist.YearToDate)
	pf.Logger.Infof("Prefetching data for %s to %s...", begin, end)

	// Use a derived context to allow for cancellation if needed
	prefetchCtx, cancel := context.WithTimeout(ctx, 23*time.Hour) // Slightly less than 24h to allow shutdown
	defer cancel()

	// Attempt to fetch all data
	_, err := pf.OrchestrateService.GetAllData(prefetchCtx, persist.CorporationIDs, persist.AllianceIDs, persist.CharacterIDs, begin, end)
	if err != nil {
		pf.Logger.Infof("Error fetching detailed killmails: %v", err)
	}
}

// Stop gracefully stops the prefetching process.
// It waits for the current prefetch to complete.
func (pf *PrefetchService) Stop() {
	pf.runningMu.Lock()
	defer pf.runningMu.Unlock()

	if !pf.running {
		pf.Logger.Info("PrefetchService is not running. Stop request ignored.")
		return
	}

	// Signal the run loop to stop
	close(pf.stopChan)

	// Wait for the run loop to acknowledge and stop
	<-pf.stopped

	// Wait for all prefetch operations to finish
	pf.wg.Wait()

	// Reset channels for potential future restarts
	pf.stopChan = make(chan struct{})
	pf.stopped = make(chan struct{})
}
