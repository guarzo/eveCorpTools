package persist

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
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

// LoadLootSplits reads the loot splits from the specified JSON file.
// It returns an empty slice if the file does not exist or is empty.
func LoadLootSplits(filename string) ([]model.LootSplit, error) {
	var lootSplits []model.LootSplit

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist; return empty slice
			return lootSplits, nil
		}
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Check if file is empty
	fi, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", filename, err)
	}

	if fi.Size() == 0 {
		// Empty file; return empty slice
		return lootSplits, nil
	}

	// Reset file pointer to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file %s: %w", filename, err)
	}

	// Attempt to decode JSON
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&lootSplits); err != nil {
		if err == io.EOF {
			// Empty JSON; return empty slice
			return lootSplits, nil
		}
		return nil, fmt.Errorf("failed to decode JSON from file %s: %w", filename, err)
	}

	return lootSplits, nil
}

// SaveLootSplits writes the provided loot splits to the specified JSON file.
func SaveLootSplits(filename string, lootSplits []model.LootSplit) error {
	return WriteJSONToFile(filename, lootSplits)
}

// AddLootSplit adds a new loot split to the existing splits and saves them.
// It sets the Date field of the new split to the current UTC time.
func AddLootSplit(filename string, newSplit model.LootSplit) error {
	// Load existing splits
	lootSplits, err := LoadLootSplits(filename)
	if err != nil {
		return fmt.Errorf("failed to load existing loot splits: %w", err)
	}

	// Set the Date field to current UTC time
	newSplit.Date = time.Now().UTC().Format(time.RFC3339)

	// Append the new split
	lootSplits = append(lootSplits, newSplit)

	// Save all splits
	if err := SaveLootSplits(filename, lootSplits); err != nil {
		return fmt.Errorf("failed to save loot splits: %w", err)
	}

	return nil
}

// DeleteLootSplit deletes a loot split by its ID and saves the updated splits.
func DeleteLootSplit(filename string, id int) ([]model.LootSplit, error) {
	// Load existing splits
	lootSplits, err := LoadLootSplits(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load existing loot splits: %w", err)
	}

	// Validate ID
	if id < 0 || id >= len(lootSplits) {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	// Remove the split with the given ID
	lootSplits = append(lootSplits[:id], lootSplits[id+1:]...)

	// Save the updated splits
	if err := SaveLootSplits(filename, lootSplits); err != nil {
		return nil, fmt.Errorf("failed to save updated loot splits: %w", err)
	}

	return lootSplits, nil
}

// CreateLootSplitBackup creates a backup of the loot splits file with a timestamp
func CreateLootSplitBackup() error {
	// Ensure backup directory exists
	if err := os.MkdirAll(config.LootDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("%s/loot_splits_backup_%s.json", config.LootDir, timestamp)

	// Load the current loot splits to create a backup
	lootSplits, err := LoadLootSplits(config.LootFile)
	if err != nil {
		return fmt.Errorf("failed to load current loot splits for backup: %v", err)
	}

	// Use writeJSONToFile helper to create the backup file
	return WriteJSONToFile(backupFilename, lootSplits)
}
