package fetch

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gambtho/zkillanalytics/internal/model"
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

	_, err := FetchAllData(client, persist.CorporationIDs, persist.AllianceIDs, persist.CharacterIDs, begin, end)
	if err != nil {
		fmt.Println("Error fetching detailed killmails:", err)
	}

	// Schedule the next run of prefetch in 24 hours
	time.AfterFunc(24*time.Hour, prefetch)
}

func RefreshEsiData(chartData *model.ChartData, client *http.Client) {
	fmt.Println("Refreshing ESI data...")

	emptyESI := false

	if len(chartData.ESIData.CharacterInfos) == 0 {
		fmt.Println("Empty ESI file provided to refresh")
		emptyESI = true
	}

	newEsiData := &model.ESIData{
		AllianceInfos:    make(map[int]model.Alliance),
		CharacterInfos:   make(map[int]model.Character),
		CorporationInfos: make(map[int]model.Corporation),
	}

	fmt.Println(fmt.Sprintf("Refreshing ESI data for characters... %d killmails to process", len(chartData.KillMails)))
	for index, detailedKillMail := range chartData.KillMails {
		if index%100 == 0 {
			fmt.Println(fmt.Sprintf("Processing killmail %d...%d of %d", detailedKillMail.EsiKillMail.KillMailID, index, len(chartData.KillMails)))
		}
		err := fetchAllESI(client, &detailedKillMail.EsiKillMail, newEsiData)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to fetch ESI data %s", err))
			if !emptyESI {
				fmt.Println("Error fetching ESI data, using existing data")
				return
			}
		}
	}

	chartData.ESIData = *newEsiData
	fmt.Println("ESI data refreshed.")
}

func init() {
	go prefetch()
}
