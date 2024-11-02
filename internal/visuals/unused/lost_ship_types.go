package unused

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

type ShipLossData struct {
	Name      string
	LossCount int
}

func GetLostShipTypes(chartData *model.ChartData, os *service.OrchestrateService) []ShipLossData {
	shipLosses := make(map[string]int)

	for _, km := range chartData.KillMails {
		victim := km.EsiKillMail.Victim

		if config.TrackedCharacterID(victim.CharacterID) || config.TrackedCorporationID(victim.CorporationID) {
			shipName := os.LookupType(victim.ShipTypeID)
			if shipName == "" || shipName == "Capsule" {
				continue
			}
			shipLosses[shipName]++
		}
	}

	var shipData []ShipLossData
	for name, count := range shipLosses {
		shipData = append(shipData, ShipLossData{
			Name:      name,
			LossCount: count,
		})
	}

	sort.Slice(shipData, func(i, j int) bool {
		return shipData[i].LossCount > shipData[j].LossCount
	})

	return shipData
}
