package loot

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/guarzo/zkillanalytics/internal/handlers"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

func LootAppraisalPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/tmpl/lootappraisal.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SaveLootSplitHandler handles saving the loot split details
func SaveLootSplitHandler(w http.ResponseWriter, r *http.Request) {
	var lootSplit model.LootSplit
	if err := json.NewDecoder(r.Body).Decode(&lootSplit); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lootSplit.Date = time.Now().UTC().Format(time.RFC3339)

	var lootSplits []model.LootSplit

	// Load existing splits
	file, err := os.Open("data/loot_split.json")
	if err == nil {
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&lootSplits); err != nil {
			log.Printf("Error decoding existing splits: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Append new split
	lootSplits = append(lootSplits, lootSplit)

	// Save all splits
	file, err = os.Create("data/loot_split.json")
	if err != nil {
		log.Printf("Error creating file for saving splits: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(lootSplits); err != nil {
		log.Printf("Error encoding splits to file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SaveLootSplitsHandler handles saving the loot split details
func SaveLootSplitsHandler(w http.ResponseWriter, r *http.Request) {
	var lootSplits []model.LootSplit
	if err := json.NewDecoder(r.Body).Decode(&lootSplits); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save all splits
	file, err := os.Create("data/loot_split.json")
	if err != nil {
		log.Printf("Error creating file for saving splits: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(lootSplits); err != nil {
		log.Printf("Error encoding splits to file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func FetchLootSplitsHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data/loot_split.json")
	if err != nil {
		log.Printf("Error opening file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var lootSplits []model.LootSplit
	if err := json.NewDecoder(file).Decode(&lootSplits); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(lootSplits); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LootSummaryHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/tmpl/lootsummary.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// DeleteLootSplitHandler handles the deletion of a loot split
func DeleteLootSplitHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var lootSplits []model.LootSplit

	// Load existing splits
	file, err := os.Open("data/loot_split.json")
	if err == nil {
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&lootSplits); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Remove the split with the given ID
	if requestData.ID >= 0 && requestData.ID < len(lootSplits) {
		lootSplits = append(lootSplits[:requestData.ID], lootSplits[requestData.ID+1:]...)
	} else {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Save the updated splits
	file, err = os.Create("data/loot_split.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(lootSplits); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
