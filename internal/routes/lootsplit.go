package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gambtho/zkillanalytics/internal/api/esi"
	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

func LootAppraisalPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/lootappraisal.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func FetchCharacterNamesHandler(w http.ResponseWriter, r *http.Request) {
	characterIDs := persist.CharacterIDs
	client := &http.Client{}
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	characterNames := make([]string, 0, len(characterIDs))

	for _, id := range characterIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			character, err := esi.FetchCharacterInfo(client, id)
			if err != nil {
				log.Printf("Error fetching character info for ID %d: %v", id, err)
				return
			}
			mu.Lock()
			characterNames = append(characterNames, character.Name)
			mu.Unlock()
		}(id)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(characterNames); err != nil {
		log.Printf("Error encoding character names: %v", err)
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
	tmpl, err := template.ParseFiles("static/lootsummary.tmpl")
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
