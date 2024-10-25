package cmd

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/routes"
	"github.com/gorilla/mux"
)

// logRequestHost middleware logs the host of each request
func logRequestHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("Host: %s\n", r.Host)
		next.ServeHTTP(w, r)
	})
}

// hostBasedRouting middleware routes requests based on the host
func hostBasedRouting(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("Routing based on host: %s\n", r.Host)
		switch r.Host {
		case "loot.zoolanders.space":
			//fmt.Println("Serving loot routes")
			lootRouter := mux.NewRouter()
			registerLootRoutes(lootRouter)
			lootRouter.ServeHTTP(w, r)
		case "tps.zoolanders.space":
			//fmt.Println("Serving tps routes")
			tpsRouter := mux.NewRouter()
			registerTPSRoutes(tpsRouter)
			tpsRouter.ServeHTTP(w, r)
		default:
			//fmt.Println("Serving default route")
			next.ServeHTTP(w, r)
		}
	})
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

// registerTPSRoutes registers the routes for the TPS subdomain
func registerTPSRoutes(r *mux.Router) {
	r.HandleFunc("/", routes.ServeRoute(persist.Snippets)).Methods("GET")
	r.HandleFunc("/lastMonth", routes.ServeRoute(persist.All)).Methods("GET")
	r.HandleFunc("/currentMonth", routes.ServeRoute(persist.All)).Methods("GET")
	r.HandleFunc("/config", routes.ServeRoute(persist.Config)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.NotFoundHandler = http.HandlerFunc(routes.NotFoundHandler)
}

// StartServer starts the HTTP server with the specified routes
func StartServer(port int) {
	fmt.Printf("------------------------------\n")
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("GOMEMLIMIT: %d\n", debug.SetMemoryLimit(-1))
	fmt.Printf("Version: %s\n", os.Getenv("VERSION"))
	fmt.Printf("------------------------------\n")

	mainRouter := mux.NewRouter()

	mainRouter.Use(logRequestHost)
	mainRouter.Use(hostBasedRouting)

	registerTPSRoutes(mainRouter)
	registerLootRoutes(mainRouter)

	http.Handle("/", mainRouter)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on port %d...\n", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}

func init() {
	// clear cache
	_ = persist.DeleteFilesInDirectory(persist.GetChartsDirectory())

	// Start the web server
	StartServer(8080)
}
