// internal/visuals/killsByCorporation.go

package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
)

// CorporationKillCount holds the kill count data for a corporation
type CorporationKillCount struct {
	CorporationID int    `json:"corporation_id"`
	Name          string `json:"name"`
	KillCount     int    `json:"kill_count"`
}

func GetVictimsByCorp(chartData *model.ChartData) []CorporationKillCount {
	corpKillMails := make(map[int]CorporationKillCount)

	// Populate the kill count map using victims from detailed killmails
	for _, km := range chartData.KillMails {
		victimCorpID := km.EsiKillMail.Victim.CorporationID
		if persist.Contains(config.CorporationIDs, victimCorpID) {
			continue
		}
		corpInfo, exists := chartData.CorporationInfos[victimCorpID]
		if !exists || persist.Contains(config.AllianceIDs, corpInfo.AllianceID) {
			continue
		}

		if data, found := corpKillMails[victimCorpID]; found {
			data.KillCount++
			corpKillMails[victimCorpID] = data
		} else {
			corpKillMails[victimCorpID] = CorporationKillCount{
				CorporationID: victimCorpID,
				KillCount:     1,
				Name:          corpInfo.Ticker,
			}
		}
	}

	// Convert the map to a slice and sort by kill count descending
	var sortedData []CorporationKillCount
	for _, data := range corpKillMails {
		sortedData = append(sortedData, data)
	}
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].KillCount > sortedData[j].KillCount
	})

	// Limit to the top 15 corporations
	if len(sortedData) > 15 {
		sortedData = sortedData[:15]
	}

	return sortedData
}
