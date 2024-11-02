package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

type LossesData struct {
	CharacterName string
	LossesValue   float64
	LossesCount   int
	ShipType      string
	ShipCount     int
}

func GetCombinedLossData(chartData *model.ChartData) []LossesData {
	characterDataMap := make(map[string]*LossesData)
	shipLossesMap := make(map[string]int)

	for _, km := range chartData.KillMails {
		victim := km.EsiKillMail.Victim

		if config.TrackedCharacterID(victim.CharacterID) || config.TrackedCorporationID(victim.CorporationID) {
			characterInfo, exists := chartData.CharacterInfos[victim.CharacterID]
			if !exists {
				continue
			}

			characterName := characterInfo.Name
			shipName := orchestrator.LookupType(victim.ShipTypeID)

			// Update character data
			data, exists := characterDataMap[characterName]
			if !exists {
				data = &LossesData{CharacterName: characterName}
				characterDataMap[characterName] = data
			}

			data.LossesValue += km.ZKB.TotalValue
			data.LossesCount++
			data.ShipType = shipName

			// Update ship losses
			if shipName != "" && shipName != "Capsule" {
				shipLossesMap[shipName]++
			}
		}
	}

	// Convert character data map to slice
	var lossesDataSlice []LossesData
	for _, data := range characterDataMap {
		// Get ship count
		data.ShipCount = shipLossesMap[data.ShipType]
		lossesDataSlice = append(lossesDataSlice, *data)
	}

	// Sort by losses value
	sort.Slice(lossesDataSlice, func(i, j int) bool {
		return lossesDataSlice[i].LossesValue > lossesDataSlice[j].LossesValue
	})

	return lossesDataSlice
}
