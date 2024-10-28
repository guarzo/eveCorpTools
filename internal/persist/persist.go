package persist

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guarzo/zkillanalytics/internal/model"
)

const dataDirectory = "data/monthly"
const chartsDirectory = "data/charts"
const idsFile = "data/ids.json"

func GenerateRelativeDirectoryPath(dir string) string {
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, dir)
}

func GetChartsDirectory() string {
	return GenerateRelativeDirectoryPath(chartsDirectory)
}

func IntSliceToString(raw []int) string {
	var strIDs []string
	for _, id := range raw {
		strIDs = append(strIDs, strconv.Itoa(id))
	}

	return strings.Join(strIDs, ", ")
}

func HashParams(params string) string {
	hasher := md5.New()
	hasher.Write([]byte(params))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateChartFileName(dir, description, startDate, endDate, hash string) string {
	return fmt.Sprintf("%s/%s_%s_to_%s_%s.html", dir, description, startDate, endDate, hash)
}

func GenerateEsiDataFileName() string {
	return fmt.Sprintf("%s/esi-data.json", GenerateRelativeDirectoryPath(dataDirectory))
}

// GenerateZkillFileName creates a filename based on year and month.
func GenerateZkillFileName(year, month int) string {
	return fmt.Sprintf("%s/%04d-%02d-killmails.json", GenerateRelativeDirectoryPath(dataDirectory), year, month)
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

func DeleteFilesInDirectory(dir string) error {
	// Open the directory
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	// Read all file names
	files, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	// Loop over the files and remove them
	for _, file := range files {
		filePath := filepath.Join(dir, file)
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return nil
}

func SaveIdsToFile(ids *model.Ids) error {
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}

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

// Contains checks if a slice Contains a specific element
func Contains(slice []int, element int) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}
