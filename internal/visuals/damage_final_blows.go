package visuals

import (
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
)

type CharacterData struct {
	Name       string
	FinalBlows int
	DamageDone int
}

func RenderDamageDone(chartData *model.ChartData, displayType string) *charts.Bar {
	// Initialize a map to count final blows and damage done by each attacking character
	characterDataMap := make(map[string]*CharacterData)

	// Initialize a variable to track the maximum damage done
	maxDamageDone := 0

	// Populate the characterDataMap map using attackers from detailed killmails
	for _, km := range chartData.KillMails {
		for _, attacker := range km.EsiKillMail.Attackers {
			characterInfo, exists := chartData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			characterName := characterInfo.Name

			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, attacker.AllianceID) {
				// Get the character data from the map, or initialize a new one if it doesn't exist
				data, exists := characterDataMap[characterName]
				if !exists {
					data = &CharacterData{Name: characterName}
					characterDataMap[characterName] = data
				}

				// Increment the final blow count if the attacker made the final blow
				if attacker.FinalBlow {
					data.FinalBlows++
				}

				// Add the damage done by the attacker
				data.DamageDone += attacker.DamageDone

				// Update the maximum damage done if necessary
				if data.DamageDone > maxDamageDone {
					maxDamageDone = data.DamageDone
				}
			}
		}
	}

	// Convert the map to a slice of CharacterData and sort by final blow count
	var characterDataSlice []CharacterData
	for _, data := range characterDataMap {
		characterDataSlice = append(characterDataSlice, *data)
	}

	if displayType == "damage" {
		sort.Slice(characterDataSlice, func(i, j int) bool {
			return characterDataSlice[i].DamageDone > characterDataSlice[j].DamageDone
		})
	}

	if displayType == "blows" {
		sort.Slice(characterDataSlice, func(i, j int) bool {
			return characterDataSlice[i].FinalBlows > characterDataSlice[j].FinalBlows
		})
	}

	// Replace the sorted list of character names with the names from the sorted CharacterData slice
	sortedCharacters := make([]string, len(characterDataSlice))
	for i, data := range characterDataSlice {
		sortedCharacters[i] = data.Name
	}

	// Prepare data for the chart
	var finalBlowsData, damageDoneData []opts.BarData
	for i, data := range characterDataSlice {
		finalBlowsData = append(finalBlowsData, opts.BarData{Value: data.FinalBlows,
			ItemStyle: &opts.ItemStyle{
				Color: colors[i%len(colors)],
			},
		})
		// fmt.Print("Damage Done: ", data.DamageDone, " Final Blows: ", data.FinalBlows, "\n", "Name: ", data.Name, "\n")
		damageDoneData = append(damageDoneData, opts.BarData{Value: data.DamageDone,
			ItemStyle: &opts.ItemStyle{
				Color: colors[i%len(colors)],
			},
		})
	}

	// Add the appropriate series to the chart based on the displayType parameter
	if displayType == "damage" {
		bar := newBarChart("Damage Done", false)
		bar.AddSeries("Damage Done", damageDoneData)
		bar.SetXAxis(sortedCharacters)
		return bar

	} else if displayType == "blows" {
		bar := newBarChart("Final Blows", false)
		bar.AddSeries("Final Blows", finalBlowsData)
		bar.SetXAxis(sortedCharacters)
		return bar
	}

	return nil
}
