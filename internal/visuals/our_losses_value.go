package visuals

import (
	"sort"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

type CharacterValueData struct {
	Name  string
	Value float64
}

func GetOurLossesValue(chartData *model.ChartData) []CharacterValueData {
	characterValues := make(map[string]float64)

	for _, km := range chartData.KillMails {
		victim := km.EsiKillMail.Victim

		if config.TrackedCharacterID(victim.CharacterID) || config.TrackedCorporationID(victim.CorporationID) {
			characterInfo, exists := chartData.CharacterInfos[victim.CharacterID]
			if !exists {
				continue
			}

			characterName := characterInfo.Name
			characterValues[characterName] += km.ZKB.TotalValue
		}
	}

	var characterData []CharacterValueData
	for name, value := range characterValues {
		characterData = append(characterData, CharacterValueData{
			Name:  name,
			Value: value,
		})
	}

	sort.Slice(characterData, func(i, j int) bool {
		return characterData[i].Value > characterData[j].Value
	})

	return characterData
}
