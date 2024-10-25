package fetch

import (
	"github.com/gambtho/zkillanalytics/internal/persist"
	"net/http"

	"github.com/gambtho/zkillanalytics/internal/model"
)

type FetchParams struct {
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

func NewFetchParams(client *http.Client, corporations, alliances, characters []int, year int, esiData *model.ESIData, changedIDs bool, newIDs *model.Ids) FetchParams {
	return FetchParams{
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
