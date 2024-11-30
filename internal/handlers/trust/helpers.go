package trust

import (
	"html/template"
	"path/filepath"

	"github.com/guarzo/zkillanalytics/internal/handlers"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

var (
	tmpl = template.Must(template.ParseFiles(
		filepath.Join("static", "tmpl", "trustbase.tmpl"),
		filepath.Join("static", "tmpl", "trustcontent.tmpl"),
	))
)

const Title = "Who to Trust?"

func prepareHomeData(sessionValues handlers.SessionValues, identities map[int64]model.CharacterData) model.StoreData {
	trustedCharacters, err := persist.LoadTrustedCharacters()
	if err != nil {
		xlog.Logf("Error loading trusted characters %v", err)
	}

	return model.StoreData{
		Title:                 Title,
		LoggedIn:              true,
		Identities:            identities,
		TabulatorIdentities:   convertIdentitiesToTabulatorData(identities),
		MainIdentity:          sessionValues.LoggedInUser,
		TrustedCharacters:     trustedCharacters.TrustedCharacters,
		TrustedCorporations:   trustedCharacters.TrustedCorporations,
		UntrustedCharacters:   trustedCharacters.UntrustedCharacters,
		UntrustedCorporations: trustedCharacters.UntrustedCorporations,
	}
}

func isTrusted(character model.CharacterData) bool {
	trustedCharacters, _ := persist.LoadTrustedCharacters()

	for _, char := range trustedCharacters.TrustedCharacters {
		if char.CharacterID == character.CharacterID {
			return true
		}
	}
	for _, corp := range trustedCharacters.TrustedCorporations {
		if corp.CorporationID == character.CorporationID {
			return true
		}
	}
	return false
}

func convertIdentitiesToTabulatorData(identities map[int64]model.CharacterData) []map[string]interface{} {
	var tabulatorData []map[string]interface{}

	for id, characterData := range identities {
		row := map[string]interface{}{
			"CharacterID":   characterData.CharacterID,
			"CharacterName": characterData.CharacterName,
			"Portrait":      characterData.Portrait,
			"IsTrusted":     isTrusted(identities[id]),
			"CorporationID": characterData.CorporationID,
		}
		tabulatorData = append(tabulatorData, row)
	}

	return tabulatorData
}
