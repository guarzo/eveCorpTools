package service

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gambtho/zkillanalytics/internal/persist"
)

var FetchAllMutex sync.Mutex

func prefetch() {
	if !FetchAllMutex.TryLock() {
		fmt.Println("Prefetch exiting as another fetch is in progress")
		return
	}
	defer FetchAllMutex.Unlock()

	// Fetch data directly
	client := &http.Client{}
	begin, end := persist.GetDateRange(persist.YearToDate)
	fmt.Println(fmt.Sprint("Prefetching data for", begin, "to", end))

	_, err := GetAllData(client, persist.CorporationIDs, persist.AllianceIDs, persist.CharacterIDs, begin, end)
	if err != nil {
		fmt.Println("Error fetching detailed killmails:", err)
	}

	// Schedule the next run of prefetch in 24 hours
	time.AfterFunc(24*time.Hour, prefetch)
}

func init() {
	go prefetch()
}
