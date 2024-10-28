package persist

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/guarzo/zkillanalytics/internal/model"
)

func SaveIdsToFile(ids *model.Ids) error {
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	fmt.Printf("writing id files %v", ids)
	err = os.WriteFile(idsFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadIdsFromFile() (model.Ids, error) {
	var ids model.Ids

	data, err := os.ReadFile(idsFile)
	if err != nil {
		return ids, err
	}

	err = json.Unmarshal(data, &ids)
	if err != nil {
		return ids, err
	}

	fmt.Printf("ids loaded %v", ids)
	return ids, nil
}

func CheckIfIdsChanged(ids *model.Ids) (bool, *model.Ids, string) {
	// Assuming you have a way to get the old IDs
	oldIds, err := LoadIdsFromFile()
	if err != nil {
		return true, nil, err.Error()
	}

	var newCharacterIds []int
	var newCorporationIds []int
	var newAllianceIds []int
	newIDs := false

	// Check if any alliance ID has changed
	for _, id := range ids.AllianceIDs {
		if !Contains(oldIds.AllianceIDs, id) {
			newAllianceIds = append(newAllianceIds, id)
		}

	}

	// Check if any character ID has changed
	for _, id := range ids.CharacterIDs {
		if !Contains(oldIds.CharacterIDs, id) {
			newCharacterIds = append(newCharacterIds, id)
		}
	}

	// Check if any corporation ID has changed
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
