package persist

import (
	"fmt"
	"github.com/guarzo/zkillanalytics/internal/xlog"
	"os"
	"sync"

	"github.com/guarzo/zkillanalytics/internal/model"
)

const trustedCharactersFile = "data/trust/trusted_characters.json"

// Mutex for safe concurrent access
var mu sync.Mutex

// LoadTrustedCharacters loads trusted characters and corporations from a file.
func LoadTrustedCharacters() (*model.TrustedCharacters, error) {
	xlog.Logf("Loading trusted characters from file: %s", trustedCharactersFile)

	mu.Lock()
	defer mu.Unlock()

	var trustedData model.TrustedCharacters
	if err := ReadJSONFromFile(trustedCharactersFile, &trustedData); err != nil {
		if os.IsNotExist(err) {
			xlog.Logf("Trusted characters file not found. Initializing empty trusted data.")
			return &model.TrustedCharacters{
				TrustedCharacters:     []model.TrustedCharacter{},
				TrustedCorporations:   []model.TrustedCorporation{},
				UntrustedCharacters:   []model.TrustedCharacter{},
				UntrustedCorporations: []model.TrustedCorporation{},
			}, nil
		}
		xlog.Logf("Error reading trusted characters file: %v", err)
		return nil, fmt.Errorf("failed to open trusted characters file: %v", err)
	}

	xlog.Logf("Trusted characters successfully loaded. Counts: TrustedCharacters=%d, TrustedCorporations=%d, UntrustedCharacters=%d, UntrustedCorporations=%d",
		len(trustedData.TrustedCharacters),
		len(trustedData.TrustedCorporations),
		len(trustedData.UntrustedCharacters),
		len(trustedData.UntrustedCorporations),
	)
	return &trustedData, nil
}

// SaveTrustedCharacters saves trusted characters and corporations to a file.
func SaveTrustedCharacters(trustedData *model.TrustedCharacters) error {
	mu.Lock()
	defer mu.Unlock()

	xlog.Logf("Saving trusted characters to file: %s", trustedCharactersFile)
	xlog.Logf("Counts: TrustedCharacters=%d, TrustedCorporations=%d, UntrustedCharacters=%d, UntrustedCorporations=%d",
		len(trustedData.TrustedCharacters),
		len(trustedData.TrustedCorporations),
		len(trustedData.UntrustedCharacters),
		len(trustedData.UntrustedCorporations),
	)

	// Log the first 5 character names (for debugging) without overwhelming logs
	for i, char := range trustedData.TrustedCharacters {
		if i >= 5 {
			xlog.Logf("...and more characters not shown")
			break
		}
		xlog.Logf("CharacterID: %d, CharacterName: %s", char.CharacterID, char.CharacterName)
	}

	return WriteJSONToFile(trustedCharactersFile, trustedData)
}
