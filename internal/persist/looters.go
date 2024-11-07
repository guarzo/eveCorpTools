package persist

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	addedPilotsFile   = "data/loot/added_pilots.json"
	removedPilotsFile = "data/loot/removed_pilots.json"
	trustMu           sync.Mutex
)

// LoadAddedPilots loads the list of added pilot names.
func LoadAddedPilots() (map[string]bool, error) {
	trustMu.Lock()
	defer trustMu.Unlock()

	file, err := os.Open(addedPilotsFile)
	if os.IsNotExist(err) {
		return make(map[string]bool), nil // No pilots added yet
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pilots map[string]bool
	if err := json.NewDecoder(file).Decode(&pilots); err != nil {
		return nil, err
	}
	return pilots, nil
}

// AddPilotName adds a pilot name to the added pilots JSON file.
func AddPilotName(name string) error {
	trustMu.Lock()
	defer trustMu.Unlock() // Ensure mutex is always released

	// Load existing pilots or create an empty map if the file does not exist.
	pilots := make(map[string]bool)
	if err := ReadJSONFromFile(addedPilotsFile, &pilots); err != nil {
		if os.IsNotExist(err) {
			log.Printf("File %s does not exist, initializing with an empty map.", addedPilotsFile)
		} else {
			return fmt.Errorf("error reading added pilots file: %v", err)
		}
	}

	log.Printf("Adding pilot %s to the pilots map", name)
	pilots[name] = true // Mark the pilot as added

	// Save updated pilots
	if err := WriteJSONToFile(addedPilotsFile, pilots); err != nil {
		log.Printf("Error saving added pilots to file %s: %v", addedPilotsFile, err)
		return err
	}

	log.Println("Pilot successfully added and saved to file")
	return nil
}

// LoadRemovedPilots loads the list of removed pilot names.
func LoadRemovedPilots() (map[string]bool, error) {
	trustMu.Lock()
	defer trustMu.Unlock()

	file, err := os.Open(removedPilotsFile)
	if os.IsNotExist(err) {
		return make(map[string]bool), nil // No pilots removed yet
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pilots map[string]bool
	if err := json.NewDecoder(file).Decode(&pilots); err != nil {
		return nil, err
	}
	return pilots, nil
}

// RemovePilot marks a pilot as removed by adding their name to the removed pilots JSON file.
func RemovePilot(name string) error {
	trustMu.Lock()
	defer trustMu.Unlock() // Ensure mutex is always released

	// Load existing removed pilots or create an empty map if the file does not exist.
	pilots := make(map[string]bool)
	if err := ReadJSONFromFile(removedPilotsFile, &pilots); err != nil {
		if os.IsNotExist(err) {
			log.Printf("File %s does not exist, initializing with an empty map.", removedPilotsFile)
		} else {
			return fmt.Errorf("error reading removed pilots file: %v", err)
		}
	}

	log.Printf("Marking pilot %s as removed", name)
	pilots[name] = true

	// Save updated removed pilots
	if err := WriteJSONToFile(removedPilotsFile, pilots); err != nil {
		log.Printf("Error saving removed pilots to file %s: %v", removedPilotsFile, err)
		return err
	}

	log.Println("Pilot successfully marked as removed and saved to file")
	return nil
}
