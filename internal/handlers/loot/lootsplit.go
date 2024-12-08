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

	// Load existing loot splits to determine the next available ID
	existingSplits, err := persist.LoadLootSplits(config.LootFile)
	if err != nil {
		log.Printf("Error loading existing loot splits: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	nextID := 1
	if len(existingSplits) > 0 {
		nextID = existingSplits[len(existingSplits)-1].ID + 1
	}

	lootSplit.ID = nextID

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

	// Load existing loot splits to determine the next available ID
	existingSplits, err := persist.LoadLootSplits(config.LootFile)
	if err != nil {
		log.Printf("Error loading existing loot splits: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	nextID := 1
	if len(existingSplits) > 0 {
		nextID = existingSplits[len(existingSplits)-1].ID + 1
	}

	// Assign IDs to new splits
	for i := range lootSplits {
		lootSplits[i].ID = nextID
		nextID++
	}

	// Combine old and new splits
	combinedSplits := append(existingSplits, lootSplits...)

	// Save all splits back to the file
	if err := persist.SaveLootSplits(config.LootFile, combinedSplits); err != nil {
		log.Printf("Error saving loot splits: %v", err)
		http.Error(w, "Failed to save loot splits", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FetchLootSplitsHandler handles fetching all loot splits.
func FetchLootSplitsHandler(w http.ResponseWriter, r *http.Request) {
	_ = BackfillIDs()
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

func DeleteLootSplitHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	lootSplits, err := persist.LoadLootSplits(config.LootFile)
	if err != nil {
		log.Printf("Error loading loot splits: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Filter out the record to delete
	updatedSplits := make([]model.LootSplit, 0)
	deleted := false
	for _, split := range lootSplits {
		if split.ID != requestData.ID {
			updatedSplits = append(updatedSplits, split)
		} else {
			deleted = true
		}
	}

	if !deleted {
		log.Printf("No loot split found with ID: %d", requestData.ID)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := persist.SaveLootSplits(config.LootFile, updatedSplits); err != nil {
		log.Printf("Error saving updated loot splits: %v", err)
		http.Error(w, "Failed to save updated loot splits", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted loot split with ID: %d", requestData.ID)
	w.WriteHeader(http.StatusOK)
}

func UpdateLootSplitHandler(w http.ResponseWriter, r *http.Request) {
	var updateRequest struct {
		ID           int    `json:"id"`
		BattleReport string `json:"battleReport"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		log.Printf("Error decoding update request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	lootSplits, err := persist.LoadLootSplits(config.LootFile)
	if err != nil {
		log.Printf("Error loading loot splits: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	updated := false
	for i, split := range lootSplits {
		if split.ID == updateRequest.ID {
			lootSplits[i].BattleReport = updateRequest.BattleReport
			updated = true
			break
		}
	}

	if !updated {
		log.Printf("No loot split found with ID: %d", updateRequest.ID)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := persist.SaveLootSplits(config.LootFile, lootSplits); err != nil {
		log.Printf("Error saving updated loot splits: %v", err)
		http.Error(w, "Failed to save updated loot split", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated BattleReport for ID: %d", updateRequest.ID)
	w.WriteHeader(http.StatusOK)
}

func BackfillIDs() error {
	lootSplits, err := persist.LoadLootSplits(config.LootFile)
	if err != nil {
		return fmt.Errorf("error loading loot splits: %v", err)
	}

	// Assign unique IDs to each loot split
	for i := range lootSplits {
		lootSplits[i].ID = i + 1 // Use a simple incrementing ID
	}

	// Save the updated loot splits back to file
	if err := persist.SaveLootSplits(config.LootFile, lootSplits); err != nil {
		return fmt.Errorf("error saving loot splits: %v", err)
	}

	return nil
}
