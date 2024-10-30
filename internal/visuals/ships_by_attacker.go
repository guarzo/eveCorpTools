package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/model"
)

type OurShipsUsedData struct {
	Characters []string         `json:"Characters"`
	ShipNames  []string         `json:"ShipNames"`
	SeriesData map[string][]int `json:"SeriesData"`
}

func GetOurShipsUsed(chartData *model.ChartData) OurShipsUsedData {
	characterShipCounts := make(map[string]map[string]int)
	shipNameSet := make(map[string]struct{})
	characters := []string{}
	shipNames := []string{}
	seriesData := make(map[string][]int)

	// Process the data
	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			if isOurCharacter(attacker.CharacterID) {
				characterInfo := chartData.CharacterInfos[attacker.CharacterID]
				characterName := characterInfo.Name

				// Assume you have a function to get ship type name by ID
				shipName := orchestrator.LookupType(attacker.ShipTypeID)

				if shipName == "" || shipName == "Capsule" || shipName == "Unknown" {
					continue
				}

				if _, exists := characterShipCounts[characterName]; !exists {
					characterShipCounts[characterName] = make(map[string]int)
				}

				characterShipCounts[characterName][shipName]++
				shipNameSet[shipName] = struct{}{}
			}
		}
	}

	// Convert sets to slices
	for shipName := range shipNameSet {
		shipNames = append(shipNames, shipName)
	}
	sort.Strings(shipNames)

	for characterName := range characterShipCounts {
		characters = append(characters, characterName)
	}
	sort.Strings(characters)

	// Prepare series data
	for _, shipName := range shipNames {
		data := make([]int, len(characters))
		for i, characterName := range characters {
			count := characterShipCounts[characterName][shipName]
			data[i] = count
		}
		seriesData[shipName] = data
	}

	return OurShipsUsedData{
		Characters: characters,
		ShipNames:  shipNames,
		SeriesData: seriesData,
	}
}
