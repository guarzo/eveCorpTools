package persist

import (
	"fmt"
	"os"

	"github.com/guarzo/zkillanalytics/internal/model"
)

const killMailDirectory = "data/tps/store"

// ReadKillMailsFromFile loads a KillMailData from a JSON file.
func ReadKillMailsFromFile(fileName string) (*model.KillMailData, error) {
	var killMailData model.KillMailData
	if err := ReadJSONFromFile(fileName, &killMailData); err != nil {
		return nil, err
	}
	return &killMailData, nil
}

// SaveKillMailsToFile saves detailed killmails to a JSON file.
func SaveKillMailsToFile(fileName string, kmData *model.KillMailData) error {
	// Ensure the directory exists
	if err := os.MkdirAll(GenerateRelativeDirectoryPath(killMailDirectory), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return WriteJSONToFile(fileName, kmData)
}
