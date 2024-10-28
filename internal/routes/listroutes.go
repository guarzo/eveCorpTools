// internal/routes/routes.go

package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// ListRoutesHandler returns a JSON list of all registered routes.
func ListRoutesHandler(router *mux.Router, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Listing all registered routes")
		var routes []string
		err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			pathTemplate, err := route.GetPathTemplate()
			if err != nil {
				pathTemplate = "unknown path"
			}
			methods, err := route.GetMethods()
			if err != nil || len(methods) == 0 {
				methods = []string{"ANY"}
			}
			routes = append(routes, fmt.Sprintf("%s [%s]", pathTemplate, strings.Join(methods, ", ")))
			return nil
		})
		if err != nil {
			http.Error(w, "Error retrieving routes", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routes)
	}
}
