package loot

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/guarzo/zkillanalytics/internal/persist"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

// LootAppraisalPageHandler renders the loot appraisal page.
func LootAppraisalPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/tmpl/lootappraisal.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// SaveLootSplitHandler handles saving a single loot split.
func SaveLootSplitHandler(w http.ResponseWriter, r *http.Request) {
	var lootSplit model.LootSplit
	if err := json.NewDecoder(r.Body).Decode(&lootSplit); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Attempt to add the new loot split
	if err := persist.AddLootSplit(config.LootFile, lootSplit); err != nil {
		log.Printf("Error saving loot split: %v", err)
		http.Error(w, "Failed to save loot split", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SaveLootSplitsHandler handles saving multiple loot splits.
func SaveLootSplitsHandler(w http.ResponseWriter, r *http.Request) {
	var lootSplits []model.LootSplit
	if err := json.NewDecoder(r.Body).Decode(&lootSplits); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Attempt to save all loot splits
	if err := persist.SaveLootSplits(config.LootFile, lootSplits); err != nil {
		log.Printf("Error saving loot splits: %v", err)
		http.Error(w, "Failed to save loot splits", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FetchLootSplitsHandler handles fetching all loot splits.
func FetchLootSplitsHandler(w http.ResponseWriter, r *http.Request) {
	lootSplits, err := persist.LoadLootSplits(config.LootFile)
	if err != nil {
		log.Printf("Error loading loot splits: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Successfully loaded splits; return them as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lootSplits); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// LootSummaryHandler renders the loot summary page.
func LootSummaryHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/tmpl/lootsummary.tmpl")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// DeleteLootSplitHandler handles the deletion of a loot split by ID.
func DeleteLootSplitHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Step 1: Create a backup of the loot splits file before deletion
	if err := persist.CreateLootSplitBackup(); err != nil {
		http.Error(w, "Failed to create backup", http.StatusInternalServerError)
		return
	}

	// Attempt to delete the loot split
	_, err := persist.DeleteLootSplit(config.LootFile, requestData.ID)
	if err != nil {
		log.Printf("Error deleting loot split: %v", err)
		if err.Error() == fmt.Sprintf("invalid ID: %d", requestData.ID) {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to delete loot split", http.StatusInternalServerError)
		}
		return
	}

	// Acknowledge successful deletion
	w.WriteHeader(http.StatusOK)
}
