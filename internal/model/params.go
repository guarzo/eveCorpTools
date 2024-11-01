package model

import (
	"net/http"
)

type Params struct {
	Client       *http.Client
	Corporations []int
	Alliances    []int
	Characters   []int
	Year         int
	EsiData      *ESIData
	ChangedIDs   bool
	NewIDs       *Ids
}

func NewParams(client *http.Client, corporations, alliances, characters []int, year int, esiData *ESIData, changedIDs bool, newIDs *Ids) Params {
	return Params{
		Client:       client,
		Corporations: corporations,
		Alliances:    alliances,
		Characters:   characters,
		Year:         year,
		EsiData:      esiData,
		ChangedIDs:   changedIDs,
		NewIDs:       newIDs,
	}
}
