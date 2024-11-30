// cmd/start.go

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/api/esi"
	"github.com/guarzo/zkillanalytics/internal/api/zkill"
	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/data"
	"github.com/guarzo/zkillanalytics/internal/handlers"
	"github.com/guarzo/zkillanalytics/internal/handlers/loot"
	"github.com/guarzo/zkillanalytics/internal/handlers/tps"
	"github.com/guarzo/zkillanalytics/internal/handlers/trust"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
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
				"method":    r.Method,
				"path":      r.URL.Path,
				"host":      r.Host,
				"remote":    r.RemoteAddr,
				"userAgent": r.UserAgent(),
			}).Info("Handling request")
			next.ServeHTTP(w, r)
		})
	}
}

// registerTPSRoutes registers the routes for the TPS subdomain
func registerTPSRoutes(r *mux.Router, orchestrateService *service.OrchestrateService, sessionStore *handlers.SessionService, esiService *service.EsiService) {
	r.Use(handlers.AuthMiddleware(sessionStore, esiService))
	r.HandleFunc("/login", handlers.LoginHandler(esiService))
	r.HandleFunc("/landing", handlers.LandingHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler(sessionStore))
	r.HandleFunc("/callback/", handlers.CallbackHandler(sessionStore, esiService))

	r.HandleFunc("/", tps.TPSHandler(config.Snippets, orchestrateService)).Methods("GET")
	r.HandleFunc("/refresh", tps.RefreshTPSHandler(orchestrateService)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
}

// registerLootRoutes registers the routes for the loot subdomain
func registerLootRoutes(r *mux.Router, sessionStore *handlers.SessionService, esiService *service.EsiService) {
	r.Use(handlers.AuthMiddleware(sessionStore, esiService))
	r.HandleFunc("/login", handlers.LoginHandler(esiService))
	r.HandleFunc("/landing", handlers.LandingHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler(sessionStore))
	r.HandleFunc("/callback/", handlers.CallbackHandler(sessionStore, esiService))

	r.HandleFunc("/", loot.LootAppraisalPageHandler).Methods("GET")
	r.HandleFunc("/loot-appraisal", loot.LootAppraisalPageHandler).Methods("GET")
	r.HandleFunc("/appraise-loot", loot.AppraiseLootHandler).Methods("POST")
	r.HandleFunc("/save-loot-split", loot.SaveLootSplitHandler).Methods("POST")
	r.HandleFunc("/delete-loot-split", loot.DeleteLootSplitHandler).Methods("POST")
	r.HandleFunc("/save-loot-splits", loot.SaveLootSplitsHandler).Methods("POST")
	r.HandleFunc("/fetch-loot-splits", loot.FetchLootSplitsHandler).Methods("GET")
	r.HandleFunc("/update-loot-split", loot.UpdateLootSplitHandler).Methods("POST")

	r.HandleFunc("/loot-summary", loot.LootSummaryHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
}

// registerTrustRoutes registers the routes for the loot subdomain
func registerTrustRoutes(r *mux.Router, sessionStore *handlers.SessionService, trustedService *service.TrustedService, esiService *service.EsiService) {
	r.Use(handlers.AuthMiddleware(sessionStore, esiService))
	r.HandleFunc("/login", handlers.LoginHandler(esiService))
	r.HandleFunc("/landing", handlers.LandingHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler(sessionStore))
	r.HandleFunc("/callback/", handlers.CallbackHandler(sessionStore, esiService))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// user functions
	r.HandleFunc("/auth-character", handlers.AuthCharacterHandler(esiService))
	r.HandleFunc("/", trust.HomeHandler(sessionStore, esiService))

	r.HandleFunc("/update-comment", trust.UpdateCommentHandler)
	r.HandleFunc("/update-is-on-couch", trust.UpdateIsOnCouchHandler)

	r.HandleFunc("/validate-and-add-trusted-character", trust.AddTrustedCharacterHandler(sessionStore, trustedService, esiService)).Methods("POST")
	r.HandleFunc("/remove-trusted-character", trust.RemoveTrustedCharacterHandler(trustedService)).Methods("POST")

	r.HandleFunc("/validate-and-add-trusted-corporation", trust.AddTrustedCorporationHandler(sessionStore, trustedService, esiService)).Methods("POST")
	r.HandleFunc("/remove-trusted-corporation", trust.RemoveTrustedCorporationHandler(trustedService)).Methods("POST")

	r.HandleFunc("/add-contacts", trust.AddContactsHandler(sessionStore, esiService))
	r.HandleFunc("/delete-contacts", trust.DeleteContactsHandler(sessionStore, esiService))

	r.HandleFunc("/validate-and-add-untrusted-character", trust.AddUntrustedCharacterHandler(sessionStore, trustedService, esiService)).Methods("POST")
	r.HandleFunc("/remove-untrusted-character", trust.RemoveUntrustedCharacterHandler(trustedService)).Methods("POST")

	r.HandleFunc("/validate-and-add-untrusted-corporation", trust.AddUntrustedCorporationHandler(sessionStore, trustedService, esiService)).Methods("POST")
	r.HandleFunc("/remove-untrusted-corporation", trust.RemoveUntrustedCorporationHandler(trustedService)).Methods("POST")

	// admin routes
	r.HandleFunc("/reset-identities", handlers.ResetIdentitiesHandler(sessionStore))
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

}

// registerDefaultRoutes registers the default routes for hosts like localhost:8080 or zoolanders.space
func registerDefaultRoutes(r *mux.Router, logger *logrus.Logger) {
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
	r.HandleFunc("/routes", handlers.ListRoutesHandler(r, logger)).Methods("GET")

	// Add a handler for the root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Default Router!"))
	}).Methods("GET")

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}).Methods("GET")
}

// StartServer starts the HTTP server with the specified routes
func StartServer(setup *config.AppSetup) {
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
		"VERSION":        setup.Version,
		"Listen Address": fmt.Sprintf(":%d", setup.Port),
		"HOST_CONFIG":    setup.HostConfig,
	}).Info("Runtime information")

	logger.Infof("host is %s", setup.HostConfig)

	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		logger.Errorf("error attaching mime extension %v]", err)
	}

	// Validate host_config
	validHosts := map[string]bool{
		"tps.zoolanders.space":   true,
		"loot.zoolanders.space":  true,
		"trust.zoolanders.space": true,
		"localhost":              true, // Optionally include "localhost" as a valid default
	}

	if setup.HostConfig != "" && !validHosts[setup.HostConfig] {
		logger.Fatalf("Invalid host_config: %s. Must be one of %v", setup.HostConfig, keys(validHosts))
	}

	// Initialize Cache
	cache := persist.NewCache(logger)
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

	failedChars, err := persist.LoadFailedCharacters()
	if err != nil {
		logger.Errorf("Failed to load failed characters: %v", err)
	}

	// Initialize configuration directory
	if err = persist.Initialize(setup.Key); err != nil {
		log.Fatalf("Failed to initialize identity: %v", err)
	}

	// Ensure necessary directories exist
	dirs := []string{
		"data",
		"data/tps",
		"data/tps/store",
		"data/tps/charts",
		"data/loot",
		"data/trust",
	}
	for _, dir := range dirs {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			logger.Fatalf("Failed to create directory %s: %v", dir, err)
		} else {
			logger.Infof("Ensured directory exists: %s", dir)
		}
	}

	// Clear cache directory if needed
	err = persist.DeleteFilesInDirectory(persist.GetChartsDirectory())
	if err != nil {
		logger.Infof("Using new charts directory: %v", err)
	}

	// Initialize HTTP Client with User-Agent
	httpClient := utils.NewHTTPClientWithUserAgent(setup.UserAgent)

	lootSessionStore, lootEsiService := initializeForHost("loot", failedChars, httpClient, cache, logger, setup.Secret)
	tpsSessionStore, tpsEsiService := initializeForHost("tps", failedChars, httpClient, cache, logger, setup.Secret)
	trustSessionStore, trustEsiService := initializeForHost("trust", failedChars, httpClient, cache, logger, setup.Secret)

	zkillClient := zkill.NewZkillClient(config.ZkillURL, httpClient, cache, logger)
	invTypeService := data.NewInvTypeService(logger) // Ensure this function exists and is correctly implemented
	err = invTypeService.LoadInvTypes()
	if err != nil {
		logger.Fatalf("failed to load invtypes %v", err)
	}
	killMailService := service.NewKillMailService(zkillClient, tpsEsiService, cache, logger)
	orchestrateService := service.NewOrchestrateService(tpsEsiService, killMailService, invTypeService, failedChars, cache, logger, httpClient)
	// Load trusted characters on startup
	dataLoader := persist.LoadTrustedCharacters
	dataSaver := persist.SaveTrustedCharacters

	// Initialize TrustedService with dependency injection
	trustedService := service.NewTrustedService(dataLoader, dataSaver, logger)

	// Create a root context that we can cancel on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure resources are cleaned up

	// Initialize and start PrefetchService with the root context
	prefetchService := service.NewPrefetchService(orchestrateService, logger)
	prefetchService.Start(ctx)

	// Initialize Main Router
	mainRouter := mux.NewRouter()

	// Apply Middlewares
	mainRouter.Use(loggingMiddleware(logger)) // Detailed request logging
	mainRouter.Use(logRequestHost(logger))    // Existing host/path logging

	// Function to create a host matcher that handles hostConfig
	hostMatcher := func(targetHost string) mux.MatcherFunc {
		return func(r *http.Request, rm *mux.RouteMatch) bool {
			host := utils.GetHost(r.Host)

			// Remove port if present
			if idx := strings.Index(host, ":"); idx != -1 {
				host = host[:idx]
			}

			// Allow 'localhost' to match any targetHost for development purposes
			isLocalhost := strings.EqualFold(host, "localhost")

			//logger.Infof("host is %s", host)
			//logger.Infof("hostConfig is %s", setup.HostConfig)

			var match bool
			if isLocalhost && setup.HostConfig == targetHost {
				match = true
			} else {
				match = strings.EqualFold(host, targetHost)
			}

			//logger.WithFields(logrus.Fields{
			//	"originalHost":  r.Host,
			//	"effectiveHost": host,
			//	"targetHost":    targetHost,
			//	"match":         match,
			//}).Info("Host matching")
			return match
		}
	}

	// Initialize Subrouters with Host Matchers
	tpsRouter := mainRouter.MatcherFunc(hostMatcher("tps.zoolanders.space")).Subrouter()
	registerTPSRoutes(tpsRouter, orchestrateService, tpsSessionStore, tpsEsiService)
	logger.Info("Registered TPS subdomain routes")

	lootRouter := mainRouter.MatcherFunc(hostMatcher("loot.zoolanders.space")).Subrouter()
	registerLootRoutes(lootRouter, lootSessionStore, lootEsiService)
	logger.Info("Registered Loot subdomain routes")

	trustRouter := mainRouter.MatcherFunc(hostMatcher("trust.zoolanders.space")).Subrouter()
	registerTrustRoutes(trustRouter, trustSessionStore, trustedService, trustEsiService)
	logger.Info("Registered Trust subdomain routes")

	// Default Router handles all other hosts
	defaultRouter := mainRouter.NewRoute().Subrouter()
	registerDefaultRoutes(defaultRouter, logger)
	logger.Info("Registered Default routes")

	// Log all registered routes for debugging
	handlers.ListRoutes(mainRouter, logger)

	// Implement a catch-all NotFoundHandler
	mainRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"host":   r.Host,
		}).Warn("Route not found")
		http.Error(w, "404 page not found", http.StatusNotFound)
	})

	// Define server address
	addr := fmt.Sprintf(":%d", setup.Port)
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
		logger.Infof("Starting server version %s on port %d", setup.Version, setup.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			// Error starting or closing listener
			logger.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// Block until graceful shutdown is complete
	<-idleConnsClosed
	logger.Info("Server gracefully stopped")
}

func initializeForHost(host string, failedChars *model.FailedCharacters, httpClient *http.Client, cache *persist.Cache, logger *logrus.Logger, secret string) (*handlers.SessionService, *service.EsiService) {
	// Get environment variables for the specific host
	clientID, clientSecret, callbackURL := utils.GetESIEnv(host)

	// Initialize session store for the specific host
	sessionStore := handlers.NewSessionService(secret) // Use unique session name if needed

	// Initialize ESI client for the specific host
	esiClient := esi.NewEsiClient(config.BaseEsiURL, failedChars, httpClient, cache, logger)
	esiClient.InitializeOAuth(clientID, clientSecret, callbackURL)

	// Initialize ESI service for the specific host
	esiService := service.NewEsiService(esiClient, cache, logger)

	return sessionStore, esiService
}

// Helper function to get keys of a map
func keys(m map[string]bool) []string {
	var list []string
	for k := range m {
		list = append(list, k)
	}
	return list
}
