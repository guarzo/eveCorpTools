package loot

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/guarzo/zkillanalytics/internal/handlers"
	"github.com/guarzo/zkillanalytics/internal/persist"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
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

// GetPilotNamesHandler retrieves and returns pilot names.
func GetPilotNamesHandler(orchestrateService *service.OrchestrateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load added and removed pilots
		addedPilots, err := persist.LoadAddedPilots()
		if err != nil {
			handlers.WriteJSONError(w, "Failed to load added pilots", "", http.StatusInternalServerError, orchestrateService.Logger)
			return
		}

		removedPilots, err := persist.LoadRemovedPilots()
		if err != nil {
			handlers.WriteJSONError(w, "Failed to load removed pilots", "", http.StatusInternalServerError, orchestrateService.Logger)
			return
		}

		// Start with hardcoded list of character IDs and their names
		characterIDs := config.CharacterIDs
		var wg sync.WaitGroup
		mu := &sync.Mutex{}
		characterNames := make(map[string]bool)

		for _, id := range characterIDs {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				character, err := orchestrateService.ESIService.EsiClient.GetCharacterInfo(context.TODO(), id)
				if err != nil || character == nil {
					orchestrateService.Logger.Warnf("Error fetching character info for ID %d: %v", id, err)
					return
				}

				if _, removed := removedPilots[character.Name]; removed {
					return // Skip removed characters
				}

				mu.Lock()
				characterNames[character.Name] = true
				mu.Unlock()
			}(id)
		}

		wg.Wait()

		// Add names from added pilots if they are not in the removed list
		for name := range addedPilots {
			if _, removed := removedPilots[name]; !removed {
				characterNames[name] = true
			}
		}

		// Convert map keys to slice of names
		response := make([]string, 0, len(characterNames))
		for name := range characterNames {
			response = append(response, name)
		}

		handlers.WriteJSONResponse(w, response, http.StatusOK, orchestrateService.Logger)
	}
}

// AddPilotHandler handles adding a new pilot.
func AddPilotHandler(orchestrateService *service.OrchestrateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("AddPilotHandler invoked")

		var requestData struct {
			Name string `json:"name"`
		}

		// Decode the request body
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			log.Printf("Failed to decode request body: %v", err)
			handlers.WriteJSONError(w, "Invalid request payload", requestData.Name, http.StatusBadRequest, orchestrateService.Logger)
			return
		}

		if requestData.Name == "" {
			log.Println("Pilot name is required but missing")
			handlers.WriteJSONError(w, "Pilot name is required", requestData.Name, http.StatusBadRequest, orchestrateService.Logger)
			return
		}

		log.Printf("Attempting to add pilot: %s", requestData.Name)

		if err := persist.AddPilotName(requestData.Name); err != nil {
			log.Printf("Failed to save pilot: %v", err)
			handlers.WriteJSONError(w, "Failed to save pilot", requestData.Name, http.StatusInternalServerError, orchestrateService.Logger)
			return
		}

		// Send success response
		log.Printf("Pilot %s added successfully", requestData.Name)
		handlers.WriteJSONResponse(w, map[string]string{"status": "success"}, http.StatusOK, orchestrateService.Logger)
	}
}

// RemovePilotHandler handles removing an existing pilot.
func RemovePilotHandler(orchestrateService *service.OrchestrateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Printf("RemovePilotHandler invoked with name: %s", name)

		if name == "" {
			log.Println("Pilot name is required but missing")
			handlers.WriteJSONError(w, "Pilot name is required", name, http.StatusBadRequest, orchestrateService.Logger)
			return
		}

		// Attempt to remove pilot
		if err := persist.RemovePilot(name); err != nil {
			log.Printf("Failed to remove pilot: %v", err)
			handlers.WriteJSONError(w, "Failed to remove pilot", name, http.StatusInternalServerError, orchestrateService.Logger)
			return
		}

		// Log success and send response
		log.Printf("Pilot %s removed successfully", name)
		handlers.WriteJSONResponse(w, map[string]string{"status": "success"}, http.StatusNoContent, orchestrateService.Logger)
	}
}
