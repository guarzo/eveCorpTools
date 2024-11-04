package persist

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func DeleteCurrentMonthFile() error {
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()
	return DeleteKillMailFile(currentYear, int(currentMonth))
}

func DeleteKillMailFile(year, month int) error {
	fileName := GenerateZkillFileName(year, month)
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
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
