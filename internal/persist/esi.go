package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/zkillanalytics/internal/model"
)

const failedCharactersFile = "data/failed_characters.json"

func GenerateEsiDataFileName() string {
	return fmt.Sprintf("%s/esi-data.json", GenerateRelativeDirectoryPath(dataDirectory))
}

func ReadEsiDataFromFile(fileName string) (*model.ESIData, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var esiData *model.ESIData
	if err = json.Unmarshal(data, &esiData); err != nil {
		return nil, err
	}
	return esiData, nil
}

func SaveEsiDataToFile(fileName string, esiData *model.ESIData) error {
	// Ensure the directory exists
	if err := os.MkdirAll(GenerateRelativeDirectoryPath(dataDirectory), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	// Marshal the data into JSON
	data, err := json.MarshalIndent(esiData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %s", err)
	}

	// Write the JSON data to the file
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %s", err)
	}

	return nil
}

// LoadFailedCharacters loads the failed character IDs from file.
func LoadFailedCharacters() (*model.FailedCharacters, error) {
	var failedChars model.FailedCharacters
	data, err := os.ReadFile(failedCharactersFile)
	if err != nil {
		if os.IsNotExist(err) {
			// If file does not exist, return an empty structure
			return &model.FailedCharacters{CharacterIDs: make(map[int]bool)}, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &failedChars); err != nil {
		return nil, err
	}
	return &failedChars, nil
}

// SaveFailedCharacters saves the failed character IDs to file.
func SaveFailedCharacters(failedChars *model.FailedCharacters) error {
	data, err := json.MarshalIndent(failedChars, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal failed characters data: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(failedCharactersFile), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for failed characters file: %w", err)
	}

	if err := os.WriteFile(failedCharactersFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write failed characters file: %w", err)
	}
	return nil
}
