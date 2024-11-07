package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
)

type ShipKillData struct {
	ShipTypeID int    `json:"ShipTypeID"`
	KillCount  int    `json:"KillCount"`
	Name       string `json:"Name"`
}

func GetTopShipsKilledData(chartData *model.ChartData) []ShipKillData {
	// Initialize a map to count killmails by ship type
	shipKillCounts := make(map[int]ShipKillData)

	trackedCharacters := orchestrator.GetTrackedCharactersFromKillMails(chartData.KillMails, &chartData.ESIData)

	// Populate the kill count map using victims' ships from detailed killmails
	for _, km := range chartData.KillMails {
		victim := km.EsiKillMail.Victim
		if persist.Contains(trackedCharacters, victim.CharacterID) {
			continue
		}

		shipTypeID := km.EsiKillMail.Victim.ShipTypeID
		shipName := orchestrator.LookupType(shipTypeID) // Fetch the ship name

		if shipName == "" || shipName == "Capsule" || shipName == "#System" || shipName == "Mobile Tractor Unit" {
			continue
		}

		if data, found := shipKillCounts[shipTypeID]; found {
			data.KillCount++
			shipKillCounts[shipTypeID] = data
		} else {
			shipKillCounts[shipTypeID] = ShipKillData{
				ShipTypeID: shipTypeID,
				KillCount:  1,
				Name:       shipName,
			}
		}
	}

	// Convert the map to a slice of ShipKillData and sort by kill count
	var sortedData []ShipKillData
	for _, data := range shipKillCounts {
		sortedData = append(sortedData, data)
	}
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].KillCount > sortedData[j].KillCount
	})

	// Limit to the top 20 ships
	if len(sortedData) > 20 {
		sortedData = sortedData[:20]
	}
	//for _, data := range sortedData {
	//	fmt.Printf("Ship: %s, KillCount: %d\n", data.Name, data.KillCount)
	//}

	return sortedData
}
