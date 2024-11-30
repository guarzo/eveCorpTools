package handlers

import (
	"fmt"
	"net/http"

	"github.com/guarzo/zkillanalytics/internal/model"
)

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	RenderLandingPage(w, r)
	return
}

func RenderLandingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	title := r.URL.Query().Get("title")
	if title == "" {
		title = "Zoo Landing" // Set a default title if not provided
	}
	data := model.StoreData{Title: title}
	if err := Tmpl.ExecuteTemplate(w, "landing", data); err != nil {
		HandleErrorWithRedirect(w, r, fmt.Sprintf("Failed to render landing template: %v", err), "/")
	}
}
