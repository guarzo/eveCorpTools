package trust

import (
	"fmt"
	"net/http"

	"github.com/guarzo/zkillanalytics/internal/handlers"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

func HomeHandler(s *handlers.SessionService, esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, handlers.SessionName)
		sessionValues := handlers.GetSessionValues(session)

		if sessionValues.LoggedInUser == 0 {
			renderLandingPage(w, r)
			return
		}

		storeData, etag, canSkip := checkIfCanSkip(session, sessionValues, r)

		if canSkip {
			renderBaseTemplate(w, r, storeData)
			return
		}

		identities, err := validateIdentities(session, sessionValues, storeData, esiService)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to validate identities: %s", err.Error())
			handleErrorWithRedirect(w, r, errorMessage, "/logout")
			return
		}

		data := prepareHomeData(sessionValues, identities)

		etag, err = updateStoreAndSession(storeData, data, etag, session, r, w)
		if err != nil {
			xlog.Logf("Failed to update store and session: %v", err)
			return
		}

		renderBaseTemplate(w, r, data)
	}
}

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	renderLandingPage(w, r)
	return
}

func renderBaseTemplate(w http.ResponseWriter, r *http.Request, data model.HomeData) {
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to render base template: %v", err), "/")
	}
}

func renderLandingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data := model.HomeData{Title: Title}
	if err := tmpl.ExecuteTemplate(w, "landing", data); err != nil {
		handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to render landing template: %v", err), "/")
	}
}
