package persist

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/guarzo/zkillanalytics/internal/config"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	chartsDirectory = "data/charts"
)

// GenerateRelativeDirectoryPath creates an absolute path to a specified subdirectory.
func GenerateRelativeDirectoryPath(subDir string) string {
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, subDir)
}

// GetChartsDirectory returns the path to the charts directory.
func GetChartsDirectory() string {
	return GenerateRelativeDirectoryPath(chartsDirectory)
}

// IntSliceToString converts a slice of integers to a comma-separated string.
func IntSliceToString(raw []int) string {
	strIDs := make([]string, len(raw))
	for i, id := range raw {
		strIDs[i] = strconv.Itoa(id)
	}
	return strings.Join(strIDs, ", ")
}

// HashParams generates an MD5 hash of a given parameter string.
func HashParams(params string) string {
	hasher := md5.New()
	hasher.Write([]byte(params))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateChartFileName generates a filename for chart files with a specific hash.
func GenerateChartFileName(dir, description, startDate, endDate, hash string) string {
	return fmt.Sprintf("%s/%s_%s_to_%s_%s.html", dir, description, startDate, endDate, hash)
}

// DeleteFilesInDirectory removes all files in the specified directory.
func DeleteFilesInDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return nil
}

// DeleteCurrentMonthFile removes the killmail file for the current month.
func DeleteCurrentMonthFile() error {
	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())
	return DeleteKillMailFile(currentYear, currentMonth)
}

// DeleteKillMailFile deletes a killmail file for a specific year and month.
func DeleteKillMailFile(year, month int) error {
	fileName := GenerateZkillFileName(year, month)
	return os.Remove(fileName)
}

// GenerateZkillFileName creates a filename based on year and month.
func GenerateZkillFileName(year, month int) string {
	return fmt.Sprintf("%s/%04d-%02d-killmails.json", GenerateRelativeDirectoryPath(killMailDirectory), year, month)
}

// Contains checks if a slice contains a specific element.
func Contains(slice []int, element int) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// GetDateRange returns the start and end dates based on the data mode
func GetDateRange(mode config.DataMode) (startDate, endDate string) {
	currentYear := time.Now().Format("2006")
	currentMonth, _ := time.Parse("01", time.Now().Format("01"))
	log.Printf("mode is %v", mode)
	switch mode {
	case config.YearToDate:
		startDate = fmt.Sprintf("%s-01-01", currentYear)
		endDate = time.Now().Format("2006-01-02")
	case config.MonthToDate:
		startDate = fmt.Sprintf("%s-%02d-01", currentYear, currentMonth.Month())
		endDate = time.Now().Format("2006-01-02")
	case config.PreviousMonth:
		lastMonth := currentMonth.Month() - 1
		if lastMonth == 0 {
			lastMonth = 12
			currentYearInt, _ := strconv.Atoi(currentYear)
			currentYear = strconv.Itoa(currentYearInt - 1)
		}
		startDate = fmt.Sprintf("%s-%02d-01", currentYear, lastMonth)
		endDate = fmt.Sprintf("%s-%02d-%02d", currentYear, lastMonth, DaysInMonth(lastMonth, currentYear))
	default:
		fmt.Println("Invalid data mode, defaulting to year-to-date")
		startDate = fmt.Sprintf("%s-01-01", currentYear)
		endDate = time.Now().Format("2006-01-02")
	}
	fmt.Printf("Start date: %s, End date: %s, Mode %v\n", startDate, endDate, mode)
	return
}

// DaysInMonth returns the number of days in a given month and year
func DaysInMonth(month time.Month, year string) int {
	y, _ := strconv.Atoi(year)
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if (y%4 == 0 && y%100 != 0) || y%400 == 0 {
			return 29
		}
		return 28
	}
	return 0
}
