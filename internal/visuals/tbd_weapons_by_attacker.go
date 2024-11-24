package visuals

import (
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

func RenderWeaponsByCharacter(orchestrator *service.OrchestrateService, chartData *model.ChartData) *charts.Bar {
	// Initialize a map to count weapons used by each attacking character
	characterWeapons := make(map[string]map[string]int)
	characterKills := make(map[string]int)

	// Populate the characterWeapons map using attackers from detailed killmails
	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterInfo, exists := chartData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			characterName := characterInfo.Name
			weaponName := orchestrator.LookupType(attacker.WeaponTypeID)

			// If the weapon is the same as the ship, skip it
			if attacker.WeaponTypeID == attacker.ShipTypeID {
				continue
			}

			// data-clean -- for weapons we don't care about
			if weaponName == "" || weaponName == "#System" {
				// fmt.Println("skipping weapon: ", weaponName, " for character: ", characterName, " with ID: ", km.KillMail.KillMailID)
				continue
			}

			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				if _, found := characterWeapons[characterName]; !found {
					characterWeapons[characterName] = make(map[string]int)
				}
				characterWeapons[characterName][weaponName]++
				characterKills[characterName]++
			}
		}
	}

	// Convert the map to a slice of CharacterPerformanceData and sort by kill count
	var characterData []CharacterPerformanceData
	for character, kills := range characterKills {
		characterData = append(characterData, CharacterPerformanceData{
			Name:      character,
			KillCount: kills,
		})
	}
	sort.Slice(characterData, func(i, j int) bool {
		return characterData[i].KillCount > characterData[j].KillCount
	})

	// Replace the sorted list of character names with the names from the sorted CharacterPerformanceData slice
	sortedCharacters := make([]string, len(characterData))

	for i, data := range characterData {
		sortedCharacters[i] = data.Name
	}

	// Collect all weapon names and sort them
	weaponNamesMap := make(map[string]bool)
	for _, weapons := range characterWeapons {
		for weapon := range weapons {
			weaponNamesMap[weapon] = true
		}
	}
	var weaponNames []string
	for weapon := range weaponNamesMap {
		weaponNames = append(weaponNames, weapon)
	}
	sort.Strings(weaponNames)

	// Prepare data for the chart
	seriesData := make(map[string][]opts.BarData)
	for _, weapon := range weaponNames {
		seriesData[weapon] = make([]opts.BarData, len(sortedCharacters))
	}

	for i, character := range sortedCharacters {
		for _, weapon := range weaponNames {
			count := characterWeapons[character][weapon]
			seriesData[weapon][i] = opts.BarData{Value: count}
		}
	}

	return nil
}
