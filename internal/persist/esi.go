package persist

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/zkillanalytics/internal/model"
)

const (
	failedCharactersFile = "data/tps/failed_characters.json"
	idsFile              = "data/tps/ids.json"
)

// GenerateEsiDataFileName generates the filename for ESI data.
func GenerateEsiDataFileName() string {
	return fmt.Sprintf("%s/esi-data.json", GenerateRelativeDirectoryPath(killMailDirectory))
}

// ReadEsiDataFromFile reads ESI data from a JSON file into the provided structure.
func ReadEsiDataFromFile(fileName string) (*model.ESIData, error) {
	var esiData model.ESIData
	if err := ReadJSONFromFile(fileName, &esiData); err != nil {
		return nil, err
	}
	return &esiData, nil
}

// SaveEsiDataToFile saves ESI data to a JSON file.
func SaveEsiDataToFile(fileName string, esiData *model.ESIData) error {
	// Ensure the directory exists
	if err := os.MkdirAll(GenerateRelativeDirectoryPath(killMailDirectory), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return WriteJSONToFile(fileName, esiData)
}

// LoadFailedCharacters loads the failed character IDs from file.
func LoadFailedCharacters() (*model.FailedCharacters, error) {
	var failedChars model.FailedCharacters
	if err := ReadJSONFromFile(failedCharactersFile, &failedChars); err != nil {
		if os.IsNotExist(err) {
			// If file does not exist, return an empty structure
			return &model.FailedCharacters{CharacterIDs: make(map[int]bool)}, nil
		}
		return nil, err
	}
	return &failedChars, nil
}

// SaveFailedCharacters saves the failed character IDs to file.
func SaveFailedCharacters(failedChars *model.FailedCharacters) error {
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(failedCharactersFile), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for failed characters file: %v", err)
	}
	return WriteJSONToFile(failedCharactersFile, failedChars)
}

// LoadIdsFromFile loads IDs from the specified file.
func LoadIdsFromFile() (model.Ids, error) {
	var ids model.Ids
	if err := ReadJSONFromFile(idsFile, &ids); err != nil {
		return ids, err
	}
	// fmt.Printf("ids loaded %v", ids)
	return ids, nil
}

// SaveIdsToFile saves the given IDs to the specified file.
func SaveIdsToFile(ids *model.Ids) error {
	// fmt.Printf("writing id files %v", ids)
	return WriteJSONToFile(idsFile, ids)
}

// CheckIfIdsChanged compares new and old IDs to identify any differences.
func CheckIfIdsChanged(ids *model.Ids) (bool, *model.Ids, string) {
	// Load the existing IDs from file
	oldIds, err := LoadIdsFromFile()
	if err != nil {
		return true, nil, err.Error()
	}

	// Load trusted characters and corporations
	trustedCharacters, err := LoadTrustedCharacters()
	if err != nil {
		return true, nil, fmt.Sprintf("failed to load trusted characters: %v", err)
	}

	// Append trusted character IDs
	for _, char := range trustedCharacters.TrustedCharacters {
		if !Contains(ids.CharacterIDs, int(char.CharacterID)) {
			ids.CharacterIDs = append(ids.CharacterIDs, int(char.CharacterID))
		}
	}

	// Append trusted corporation IDs (if applicable to your use case)
	for _, corp := range trustedCharacters.TrustedCorporations {
		if !Contains(ids.CorporationIDs, int(corp.CorporationID)) {
			ids.CorporationIDs = append(ids.CorporationIDs, int(corp.CorporationID))
		}
	}

	var newCharacterIds, newCorporationIds, newAllianceIds []int
	newIDs := false

	// Check for new alliance IDs
	for _, id := range ids.AllianceIDs {
		if !Contains(oldIds.AllianceIDs, id) {
			newAllianceIds = append(newAllianceIds, id)
		}
	}

	// Check for new character IDs
	for _, id := range ids.CharacterIDs {
		if !Contains(oldIds.CharacterIDs, id) {
			newCharacterIds = append(newCharacterIds, id)
		}
	}

	// Check for new corporation IDs
	for _, id := range ids.CorporationIDs {
		if !Contains(oldIds.CorporationIDs, id) {
			newCorporationIds = append(newCorporationIds, id)
		}
	}

	if len(newAllianceIds) > 0 || len(newCharacterIds) > 0 || len(newCorporationIds) > 0 {
		newIDs = true
	}

	return newIDs, &model.Ids{
		AllianceIDs:    newAllianceIds,
		CharacterIDs:   newCharacterIds,
		CorporationIDs: newCorporationIds,
	}, ""
}
