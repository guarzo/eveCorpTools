package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/guarzo/zkillanalytics/internal/persist"
)

// AuthMiddleware checks if the user is logged in and has a valid session
func AuthMiddleware(sessionService *SessionService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := sessionService.Get(r, SessionName)
			sessionValues := GetSessionValues(session)

			// Check if the user is logged in
			if sessionValues.LoggedInUser == 0 {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// You can add additional checks here (e.g., ensure token exists)
			_, err := persist.GetMainIdentityToken(sessionValues.LoggedInUser)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Continue to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
