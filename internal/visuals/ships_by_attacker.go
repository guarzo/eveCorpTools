package visuals

import (
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/service"
)

// CharacterShipData holds the data for ship types used by each attacking character
type CharacterShipData struct {
	CharacterName string
	ShipName      string
	Count         int
}

func RenderOurShips(orchestrator *service.OrchestrateService, chartData *model.ChartData) *charts.Bar {
	// Initialize a map to count kills by each attacking character
	characterShips := make(map[string]map[string]int)
	characterKills := make(map[string]int)

	// Populate the characterShips map using attackers from detailed killmails
	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterInfo, exists := chartData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			characterName := characterInfo.Name
			shipName := orchestrator.LookupType(attacker.ShipTypeID)

			// data-clean -- for ships we don't care about
			if shipName == "" || shipName == "Capsule" || shipName == "#System" {
				// fmt.Println("skipping ship: ", shipName, " for character: ", characterName, " with ID: ", km.KillMail.KillMailID)
				continue
			}

			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {

				if _, found := characterShips[characterName]; !found {
					characterShips[characterName] = make(map[string]int)
				}
				characterShips[characterName][shipName]++
				characterKills[characterName]++
			}
		}
	}

	// Convert the map to a slice of CharacterKillData and sort by kill count
	var characterData []CharacterKillData
	for character, kills := range characterKills {
		characterData = append(characterData, CharacterKillData{
			Name:      character,
			KillCount: kills,
		})
	}
	sort.Slice(characterData, func(i, j int) bool {
		return characterData[i].KillCount > characterData[j].KillCount
	})

	// Replace the sorted list of character names with the names from the sorted CharacterKillData slice
	sortedCharacters := make([]string, len(characterData))

	for i, data := range characterData {
		sortedCharacters[i] = data.Name
	}

	// Collect all ship names and sort them
	shipNamesMap := make(map[string]bool)
	for _, ships := range characterShips {
		for ship := range ships {
			shipNamesMap[ship] = true
		}
	}
	var shipNames []string
	for ship := range shipNamesMap {
		shipNames = append(shipNames, ship)
	}
	sort.Strings(shipNames)

	//fmt.Println(fmt.Sprintf("sortedCharacters: %v", sortedCharacters))
	//fmt.Println(fmt.Sprintf("shipNames: %v", shipNames))
	//fmt.Println(fmt.Sprintf("count of characters: %d", len(sortedCharacters)))

	// Prepare data for the chart
	seriesData := make(map[string][]opts.BarData)
	for _, ship := range shipNames {
		seriesData[ship] = make([]opts.BarData, len(sortedCharacters))
	}

	for i, character := range sortedCharacters {
		for _, ship := range shipNames {
			count := characterShips[character][ship]
			seriesData[ship][i] = opts.BarData{Value: count}
		}
	}
	// Create a new stacked bar chart instance
	bar := newBarChart("Ship Used", true)

	bar.SetXAxis(sortedCharacters)
	for _, ship := range shipNames {
		bar.AddSeries(ship, seriesData[ship], charts.WithBarChartOpts(opts.BarChart{Stack: "total"}))
	}
	return bar
}
