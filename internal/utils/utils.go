package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gambtho/zkillanalytics/internal/persist"
)

// CheckDataAvailability checks which months within a range have data files already present.
func CheckDataAvailability(startMonth, endMonth, year int) (map[int]bool, error) {
	dataAvailability := make(map[int]bool)
	for month := startMonth; month <= endMonth; month++ {
		fileName := persist.GenerateZkillFileName(year, month)
		if fileInfo, err := os.Stat(fileName); err == nil {
			if fileInfo.Size() <= 1*1024 {
				fmt.Printf("File %s is too small (%d bytes). Marking as unavailable.\n", fileName, fileInfo.Size())
				dataAvailability[month] = false
				continue
			}
			fmt.Println(fmt.Sprintf("Data for %04d-%02d already exists.", year, month))
			dataAvailability[month] = true
			if month == int(time.Now().Month()) {
				fileDate := fileInfo.ModTime().Truncate(24 * time.Hour)
				today := time.Now().Truncate(24 * time.Hour)
				if fileDate.Before(today) {
					fmt.Println(fmt.Sprintf("Removing stale month to date file %s...", fileName))
					dataAvailability[month] = false
					err = os.Remove(fileName)
					if err != nil {
						fmt.Println(fmt.Sprintf("Error removing stale month to date file %s: %s", fileName, err))
					}
				} else {
					fmt.Println(fmt.Sprintf("Continuing to use current month data for %s, %s", fileName, fileDate.Format("2006-01-02:15:04:05")))
				}
			}
		} else {
			dataAvailability[month] = false
		}
	}
	return dataAvailability, nil
}