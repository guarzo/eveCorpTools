package visuals

import (
	"github.com/guarzo/zkillanalytics/internal/model"
)

type SunburstData struct {
	Name     string         `json:"name"`
	Value    int            `json:"value,omitempty"`
	Children []SunburstData `json:"children,omitempty"`
}

func GetVictimsSunburst(chartData *model.ChartData) []SunburstData {
	allianceMap := make(map[int]*SunburstData)

	//for _, km := range chartData.KillMails {
	//	victim := km.EsiKillMail.Victim
	//	corpID := victim.CorporationID
	//
	//
	//	if persist.Contains(config.CorporationIDs, corpID)  {
	//		continue
	//	}
	//
	//	// Get alliance data
	//	allianceData, exists := chartData.AllianceInfos[allianceID]
	//	if !exists {
	//		allianceName := chartData.AllianceInfos[allianceID].Name
	//		allianceData = &SunburstData{Name: allianceName}
	//		allianceMap[allianceID] = allianceData
	//	}
	//
	//	// Get corporation data
	//	corpName := chartData.CorporationInfos[corpID].Name
	//	var corpData *SunburstData
	//	for i := range allianceData.Children {
	//		if allianceData.Children[i].Name == corpName {
	//			corpData = &allianceData.Children[i]
	//			break
	//		}
	//	}
	//	if corpData == nil {
	//		corpData = &SunburstData{Name: corpName}
	//		allianceData.Children = append(allianceData.Children, *corpData)
	//	}
	//
	//	// Increment kill count
	//	corpData.Value++
	//}

	// Convert allianceMap to a slice
	var sunburstData []SunburstData
	for _, data := range allianceMap {
		sunburstData = append(sunburstData, *data)
	}

	return sunburstData
}