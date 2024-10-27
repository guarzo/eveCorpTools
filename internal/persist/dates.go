package persist

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gambtho/zkillanalytics/internal/config"
)

// GetDateRange returns the start and end dates based on the data mode
func GetDateRange(mode config.DataMode) (startDate, endDate string) {
	currentYear := time.Now().Format("2006")
	currentMonth, _ := time.Parse("01", time.Now().Format("01"))

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
