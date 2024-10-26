package config

import (
	"net/http"

	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

type Params struct {
	Client       *http.Client
	Corporations []int
	Alliances    []int
	Characters   []int
	Year         int
	EsiData      *model.ESIData
	ChangedIDs   bool
	NewIDs       *model.Ids
	*persist.EntityLastPage
}

func NewParams(client *http.Client, corporations, alliances, characters []int, year int, esiData *model.ESIData, changedIDs bool, newIDs *model.Ids) Params {
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
