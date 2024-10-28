package persist

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/guarzo/zkillanalytics/internal/model"
)

// GenerateZkillFileName creates a filename based on year and month.
func GenerateZkillFileName(year, month int) string {
	return fmt.Sprintf("%s/%04d-%02d-killmails.json", GenerateRelativeDirectoryPath(dataDirectory), year, month)
}

// ReadKillMailDataFromFile loads a KillMailData from a JSON file.
func ReadKillMailDataFromFile(fileName string) (*model.KillMailData, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var dump *model.KillMailData
	if err = json.Unmarshal(data, &dump); err != nil {
		return nil, err
	}
	return dump, nil
}

// SaveKillMailsToFile saves detailed killmails to a JSON file
func SaveKillMailsToFile(fileName string, kmData *model.KillMailData) error {
	// Ensure the directory exists
	if err := os.MkdirAll(GenerateRelativeDirectoryPath(dataDirectory), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	// Marshal the data into JSON
	data, err := json.MarshalIndent(kmData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %s", err)
	}

	// Write the JSON data to the file
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %s", err)
	}

	return nil
}
