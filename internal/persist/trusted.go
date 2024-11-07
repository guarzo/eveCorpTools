package persist

import (
	"fmt"
	"os"
	"sync"

	"github.com/guarzo/zkillanalytics/internal/model"
)

const trustedCharactersFile = "data/trust/trusted_characters.json"

// Mutex for safe concurrent access
var mu sync.Mutex

// LoadTrustedCharacters loads trusted characters and corporations from a file.
func LoadTrustedCharacters() (*model.TrustedCharacters, error) {
	mu.Lock()
	defer mu.Unlock()

	var trustedData model.TrustedCharacters
	if err := ReadJSONFromFile(trustedCharactersFile, &trustedData); err != nil {
		if os.IsNotExist(err) {
			// If file does not exist, return an empty structure
			return &model.TrustedCharacters{
				TrustedCharacters:     []model.TrustedCharacter{},
				TrustedCorporations:   []model.TrustedCorporation{},
				UntrustedCharacters:   []model.TrustedCharacter{},
				UntrustedCorporations: []model.TrustedCorporation{},
			}, nil
		}
		return nil, fmt.Errorf("failed to open trusted characters file: %v", err)
	}
	return &trustedData, nil
}

// SaveTrustedCharacters saves trusted characters and corporations to a file.
func SaveTrustedCharacters(trustedData *model.TrustedCharacters) error {
	mu.Lock()
	defer mu.Unlock()

	// Ensure the directory exists
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for trusted characters file: %v", err)
	}
	return WriteJSONToFile(trustedCharactersFile, trustedData)
}
