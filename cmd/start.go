// cmd/start.go

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/gambtho/zkillanalytics/internal/api/esi"
	"github.com/gambtho/zkillanalytics/internal/api/zkill"
	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/data"
	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/routes"
	"github.com/gambtho/zkillanalytics/internal/service"
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
				// Serve loot routes
				lootRouter := mux.NewRouter()
				registerLootRoutes(lootRouter)
				lootRouter.ServeHTTP(w, r)
			case "tps.zoolanders.space":
				// Serve TPS routes
				tpsRouter := mux.NewRouter()
				registerTPSRoutes(tpsRouter, orchestrateService)
				tpsRouter.ServeHTTP(w, r)
			default:
				//fmt.Println("Serving default route")
				next.ServeHTTP(w, r)
			}
		})
	}
}

// registerDefaultRoutes registers the routes for the killmail subdomain or default host
func registerDefaultRoutes(r *mux.Router, esiService *service.EsiService) {
	// Serve static files if needed
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(routes.NotFoundHandler)
}

// registerTPSRoutes registers the routes for the TPS subdomain
func registerTPSRoutes(r *mux.Router, orchestrateService *service.OrchestrateService) {
	// Use ServeRoute with OrchestrateService
	r.HandleFunc("/", routes.ServeRoute(persist.Snippets, orchestrateService)).Methods("GET")
	r.HandleFunc("/lastMonth", routes.ServeRoute(persist.All, orchestrateService)).Methods("GET")
	r.HandleFunc("/currentMonth", routes.ServeRoute(persist.All, orchestrateService)).Methods("GET")
	r.HandleFunc("/config", routes.ServeRoute(persist.Config, orchestrateService)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(routes.NotFoundHandler)
}

// registerLootRoutes registers the routes for the loot subdomain
func registerLootRoutes(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/loot-appraisal", http.StatusMovedPermanently)
	}).Methods("GET")
	r.HandleFunc("/loot-appraisal", routes.LootAppraisalPageHandler).Methods("GET")
	r.HandleFunc("/appraise-loot", routes.AppraiseLootHandler).Methods("POST")
	r.HandleFunc("/fetch-character-names", routes.FetchCharacterNamesHandler).Methods("GET")
	r.HandleFunc("/save-loot-split", routes.SaveLootSplitHandler).Methods("POST")
	r.HandleFunc("/delete-loot-split", routes.DeleteLootSplitHandler).Methods("POST")
	r.HandleFunc("/save-loot-splits", routes.SaveLootSplitsHandler).Methods("POST")
	r.HandleFunc("/fetch-loot-splits", routes.FetchLootSplitsHandler).Methods("GET")
	r.HandleFunc("/loot-summary", routes.LootSummaryHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(routes.NotFoundHandler)
}

// StartServer starts the HTTP server with the specified routes
func StartServer(port int, userAgent string) {
	// Initialize Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

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

	// Clear cache directory if needed
	err := persist.DeleteFilesInDirectory(persist.GetChartsDirectory())
	if err != nil {
		logger.Errorf("Failed to clear cache directory: %v", err)
	}

	// Initialize HTTP Client with Timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

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

	// Initialize and start PrefetchService
	prefetchService := service.NewPrefetchService(orchestrateService, logger)
	prefetchService.Start(context.Background())

	// Initialize Main Router
	mainRouter := mux.NewRouter()

	// Apply Middlewares
	mainRouter.Use(logRequestHost(logger))
	mainRouter.Use(hostBasedRouting(logger, orchestrateService))

	// Register Default Routes (e.g., health check)
	mainRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

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
		<-sigint

		// Received an interrupt signal, initiate graceful shutdown
		logger.Info("Shutting down server...")

		// Create a deadline to wait for current operations to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout
			logger.Fatalf("HTTP server Shutdown: %v", err)
		}

		// Stop PrefetchService
		prefetchService.Stop()

		close(idleConnsClosed)
	}()

	// Start the server in a goroutine
	go func() {
		logger.Infof("Starting server on port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Error starting or closing listener
			logger.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// Block until graceful shutdown is complete
	<-idleConnsClosed
	logger.Info("Server gracefully stopped")
}
