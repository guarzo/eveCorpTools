// cmd/start.go

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/api/esi"
	"github.com/guarzo/zkillanalytics/internal/api/zkill"
	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/data"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/routes"
	"github.com/guarzo/zkillanalytics/internal/service"
	"github.com/guarzo/zkillanalytics/internal/utils"
)

// logRequestHost middleware logs the host and path of each incoming request
func logRequestHost(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.WithFields(logrus.Fields{
				"host": r.Host,
				"path": r.URL.Path,
			}).Info("Incoming request")
			next.ServeHTTP(w, r)
		})
	}
}

// loggingMiddleware logs detailed information about each incoming HTTP request.
func loggingMiddleware(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"host":   r.Host,
				"remote": r.RemoteAddr,
			}).Info("Handling request")
			next.ServeHTTP(w, r)
		})
	}
}

// hostBasedRouting middleware routes requests based on the host
func hostBasedRouting(logger *logrus.Logger, orchestrateService *service.OrchestrateService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.WithFields(logrus.Fields{
				"host": r.Host,
				"path": r.URL.Path,
			}).Info("Routing request based on host")

			switch r.Host {
			case "loot.zoolanders.space":
				logger.Info("Routing to Loot Router")
				lootRouter := mux.NewRouter()
				registerLootRoutes(lootRouter, orchestrateService)
				lootRouter.ServeHTTP(w, r)
			case "tps.zoolanders.space":
				logger.Info("Routing to TPS Router")
				tpsRouter := mux.NewRouter()
				registerTPSRoutes(tpsRouter, orchestrateService)
				tpsRouter.ServeHTTP(w, r)
			default:
				logger.Info("Routing to Default Router")
				next.ServeHTTP(w, r)
			}
		})
	}
}

// registerTPSRoutes registers the routes for the TPS subdomain
func registerTPSRoutes(r *mux.Router, orchestrateService *service.OrchestrateService) {
	r.HandleFunc("/", routes.ServeRoute(config.Snippets, orchestrateService)).Methods("GET")
	r.HandleFunc("/lastMonth", routes.ServeRoute(config.All, orchestrateService)).Methods("GET")
	r.HandleFunc("/currentMonth", routes.ServeRoute(config.All, orchestrateService)).Methods("GET")
	// r.HandleFunc("/config", routes.ServeRoute(persist.Config, orchestrateService)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(routes.NotFoundHandler)
}

// registerLootRoutes registers the routes for the loot subdomain
func registerLootRoutes(r *mux.Router, orchestrateService *service.OrchestrateService) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/loot-appraisal", http.StatusMovedPermanently)
	}).Methods("GET")
	r.HandleFunc("/loot-appraisal", routes.LootAppraisalPageHandler).Methods("GET")
	r.HandleFunc("/appraise-loot", routes.AppraiseLootHandler).Methods("POST")
	r.HandleFunc("/fetch-character-names", routes.FetchCharacterNamesHandler(orchestrateService)).Methods("GET")
	r.HandleFunc("/save-loot-split", routes.SaveLootSplitHandler).Methods("POST")
	r.HandleFunc("/delete-loot-split", routes.DeleteLootSplitHandler).Methods("POST")
	r.HandleFunc("/save-loot-splits", routes.SaveLootSplitsHandler).Methods("POST")
	r.HandleFunc("/fetch-loot-splits", routes.FetchLootSplitsHandler).Methods("GET")
	r.HandleFunc("/loot-summary", routes.LootSummaryHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(routes.NotFoundHandler)
}

// registerDefaultRoutes registers the default routes for hosts like localhost:8080 or zoolanders.space
func registerDefaultRoutes(r *mux.Router, orchestrateService *service.OrchestrateService, logger *logrus.Logger) {
	// Route "/" to ServeRoute with config.Snippets
	r.HandleFunc("/", routes.ServeRoute(config.Snippets, orchestrateService)).Methods("GET")

	// Route "/lastMonth" to ServeRoute with config.All
	r.HandleFunc("/lastMonth", routes.ServeRoute(config.All, orchestrateService)).Methods("GET")

	// Health Check Endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		health := struct {
			Status       string `json:"status"`
			CacheStatus  string `json:"cache_status"`
			ESIConnected bool   `json:"esi_connected"`
			// Add more fields as needed
		}{
			Status:       "OK",
			CacheStatus:  "Connected", // Implement actual cache status check
			ESIConnected: true,        // Implement actual ESI connection check
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	}).Methods("GET")

	// Register a route to list all routes (Optional)
	r.HandleFunc("/routes", routes.ListRoutesHandler(r, logger)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

// StartServer starts the HTTP server with the specified routes
func StartServer(port int, userAgent, version string) {
	// Initialize Logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel) // Set to Debug for more detailed logs

	// Choose a formatter that supports caller fields
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   false,
		DisableColors:   false,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Log runtime information
	logger.Info("Initializing server...")
	logger.WithFields(logrus.Fields{
		"GOMAXPROCS":     fmt.Sprintf("%d", runtime.GOMAXPROCS(0)),
		"GOMEMLIMIT":     "Not Set", // Adjust based on actual usage or remove if not applicable
		"VERSION":        os.Getenv("VERSION"),
		"Listen Address": fmt.Sprintf(":%d", port),
	}).Info("Runtime information")

	// Initialize Cache
	cache := persist.NewInMemoryCache(logger)
	if cache == nil {
		logger.Fatal("Failed to initialize cache")
	}

	// Load cache from file at startup
	cacheFile := persist.GenerateCacheDataFileName()
	if err := cache.LoadFromFile(cacheFile); err != nil {
		logger.Errorf("Failed to load cache from file: %v", err)
	} else {
		logger.Infof("Cache loaded from %s", cacheFile)
	}

	// Ensure necessary directories exist
	dirs := []string{
		"data",
		"data/monthly",
		"data/charts",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logger.Fatalf("Failed to create directory %s: %v", dir, err)
		} else {
			logger.Infof("Ensured directory exists: %s", dir)
		}
	}

	// Clear cache directory if needed
	err := persist.DeleteFilesInDirectory(persist.GetChartsDirectory())
	if err != nil {
		logger.Infof("Using new charts directory: %v", err)
	}

	// Initialize HTTP Client with User-Agent
	httpClient := utils.NewHTTPClientWithUserAgent(userAgent)

	esiClient := esi.NewEsiClient(config.BaseEsiURL, httpClient, cache, logger)
	zkillClient := zkill.NewZkillClient(config.ZkillURL, httpClient, cache, logger)

	// Initialize EsiService
	esiService := service.NewEsiService(esiClient, cache, logger)

	// Initialize InvTypeService
	invTypeService := data.NewInvTypeService(logger) // Ensure this function exists and is correctly implemented

	// Initialize KillMailService with EsiService
	killMailService := service.NewKillMailService(zkillClient, esiService, cache, logger)

	// Initialize OrchestrateService with EsiService and KillMailService
	orchestrateService := service.NewOrchestrateService(esiService, killMailService, invTypeService, cache, logger, httpClient)

	// Create a root context that we can cancel on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure resources are cleaned up

	// Initialize and start PrefetchService with the root context
	prefetchService := service.NewPrefetchService(orchestrateService, logger)
	prefetchService.Start(ctx)

	// Initialize Main Router
	mainRouter := mux.NewRouter()

	// Apply Middlewares
	mainRouter.Use(loggingMiddleware(logger))                    // Detailed request logging
	mainRouter.Use(logRequestHost(logger))                       // Existing host/path logging
	mainRouter.Use(hostBasedRouting(logger, orchestrateService)) // Host-based routing

	// Register Default Routes
	registerDefaultRoutes(mainRouter, orchestrateService, logger)

	// Register Subdomain Routes (handled by hostBasedRouting middleware)
	// No need to register them here; they're handled within the middleware

	// Log all registered routes for debugging
	utils.ListRoutes(mainRouter, logger)

	// Implement a catch-all NotFoundHandler (already handled within registerDefaultRoutes, but reaffirm)
	mainRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"host":   r.Host,
		}).Warn("Route not found")
		http.Error(w, "404 page not found", http.StatusNotFound)
	})

	// Define server address
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    addr,
		Handler: mainRouter,
	}

	// Channel to listen for interrupt or terminate signals
	idleConnsClosed := make(chan struct{})

	// Listen for OS signals to gracefully shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		// Catch interrupt and terminate signals
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		sig := <-sigint

		// Received an interrupt signal, initiate graceful shutdown
		logger.Infof("Received signal: %v. Shutting down server...", sig)

		// Cancel the root context to signal PrefetchService to stop
		cancel()

		// Create a deadline to wait for current operations to complete
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Attempt graceful shutdown of the HTTP server
		if err := server.Shutdown(shutdownCtx); err != nil {
			// Error from closing listeners, or context timeout
			logger.Errorf("HTTP server Shutdown: %v", err)
		}

		// Stop PrefetchService and wait for it to finish
		prefetchService.Stop()

		close(idleConnsClosed)
	}()

	// Start the server in a goroutine
	go func() {
		logger.Infof("Starting server version %s on port %d", version, port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Error starting or closing listener
			logger.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// Block until graceful shutdown is complete
	<-idleConnsClosed
	logger.Info("Server gracefully stopped")
}
