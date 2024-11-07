package trust

import (
	"fmt"
	"github.com/guarzo/zkillanalytics/internal/handlers"
	"github.com/guarzo/zkillanalytics/internal/service"
	"html"
	"html/template"
	"net/http"
	"path/filepath"
	"slices"
	"time"

	"github.com/gorilla/sessions"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

var (
	tmpl = template.Must(template.ParseFiles(
		filepath.Join("static", "tmpl", "trustbase.tmpl"),
		filepath.Join("static", "tmpl", "trustcontent.tmpl"),
		filepath.Join("static", "tmpl", "trustlanding.tmpl"),
	))
)

const Title = "Who to Trust?"

func sameIdentities(users []int64, identities map[int64]model.CharacterData) bool {
	var identitiesKeys []int64
	for k, _ := range identities {
		identitiesKeys = append(identitiesKeys, k)
	}

	if len(identities) != len(users) {
		return false
	}

	for k, _ := range identities {
		if !slices.Contains(users, k) {
			return false
		}
	}

	return true
}

func sameUserCount(session *sessions.Session, previousUsers, storeUsers int) bool {
	if previousUsers == 0 {
		return false
	}

	if previousUsers != storeUsers {
		return false
	}

	if authenticatedUsers, ok := session.Values[handlers.AllAuthenticatedCharacters].([]int64); ok {
		return previousUsers == len(authenticatedUsers)
	}

	return false
}

func handleErrorWithRedirect(w http.ResponseWriter, r *http.Request, errorMessage, redirectURL string) {
	// Set content type to HTML with UTF-8 encoding for proper character handling
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Escape the error message and redirect URL to prevent injection issues
	escapedMessage := html.EscapeString(errorMessage)
	escapedURL := html.EscapeString(redirectURL)

	// Construct HTML response with JavaScript for alert and redirection
	responseHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Error</title>
		</head>
		<body>
			<script>
				alert("%s");
				window.location.href = "%s";
			</script>
		</body>
		</html>`, escapedMessage, escapedURL)

	// Write the HTML response with embedded JavaScript
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(responseHTML))
}

func clearSession(s *handlers.SessionService, w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := s.Get(r, handlers.SessionName)
	if err != nil {
		xlog.Logf("Failed to get session to clear: %v", err)
	}

	// Clear the session
	session.Values = make(map[interface{}]interface{})

	// Save the session
	err = sessions.Save(r, w)
	if err != nil {
		xlog.Logf("Failed to save session to clear: %v", err)
	}
}

func checkIfCanSkip(session *sessions.Session, sessionValues handlers.SessionValues, r *http.Request) (model.HomeData, string, bool) {
	canSkip := true
	storeData, etag, ok := persist.Store.Get(sessionValues.LoggedInUser)
	if !ok || sessionValues.PreviousInputSubmitted == "" || sessionValues.PreviousInputSubmitted != r.FormValue("desired_destinations") {
		canSkip = false
	}
	if !sameUserCount(session, sessionValues.PreviousUserCount, len(storeData.Identities)) {
		canSkip = false
	}
	return storeData, etag, canSkip
}

// validUser now checks if a character is explicitly listed in the configuration or is a trusted character.
func validUser(character model.CharacterData) bool {
	return slices.Contains(config.CharacterIDs, int(character.CharacterID)) ||
		isTrustedCharacter(character.CharacterID)
}

// isTrustedCharacter checks if the character is in the list of trusted characters, ignoring corporations.
func isTrustedCharacter(characterID int64) bool {
	trustedCharacters, _ := persist.LoadTrustedCharacters()

	for _, char := range trustedCharacters.TrustedCharacters {
		if char.CharacterID == characterID {
			return true
		}
	}
	return false
}

func validateIdentities(session *sessions.Session, sessionValues handlers.SessionValues, storeData model.HomeData, esiService *service.EsiService) (map[int64]model.CharacterData, error) {
	identities := storeData.Identities

	authenticatedUsers, ok := session.Values[handlers.AllAuthenticatedCharacters].([]int64)
	if !ok {
		xlog.Logf("Failed to retrieve authenticated users from session")
		return nil, fmt.Errorf("failed to retrieve authenticated users from session")
	}

	needIdentityPopulation := len(authenticatedUsers) == 0 || !sameIdentities(authenticatedUsers, storeData.Identities) || time.Since(time.Unix(sessionValues.LastRefreshTime, 0)) > 15*time.Minute

	if needIdentityPopulation {
		userConfig, err := persist.LoadIdentities(sessionValues.LoggedInUser)

		if err != nil {
			xlog.Logf("Failed to load identities: %v", err)
			return nil, fmt.Errorf("failed to load identities: %w", err)
		}

		identities, err = esiService.EsiClient.PopulateIdentities(userConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to populate identities: %w", err)
		}

		if !validUser(identities[sessionValues.LoggedInUser]) {
			return nil, fmt.Errorf("not a valid user - ask in discord if you think this is a mistake")
		}

		if err = persist.SaveIdentities(sessionValues.LoggedInUser, userConfig); err != nil {
			return nil, fmt.Errorf("failed to save identities: %w", err)
		}

		session.Values[handlers.AllAuthenticatedCharacters] = getAuthenticatedCharacterIDs(identities)
		session.Values[handlers.LastRefreshTime] = time.Now().Unix()
	}

	return identities, nil
}

func getAuthenticatedCharacterIDs(identities map[int64]model.CharacterData) []int64 {
	authenticatedCharacters := make([]int64, 0, len(identities))
	for id := range identities {
		authenticatedCharacters = append(authenticatedCharacters, id)
	}
	return authenticatedCharacters
}

func prepareHomeData(sessionValues handlers.SessionValues, identities map[int64]model.CharacterData) model.HomeData {
	trustedCharacters, err := persist.LoadTrustedCharacters()
	if err != nil {
		xlog.Logf("Error loading trusted characters %v", err)
	}

	return model.HomeData{
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

func updateStoreAndSession(storeData model.HomeData, data model.HomeData, etag string, session *sessions.Session, r *http.Request, w http.ResponseWriter) (string, error) {
	newEtag, err := persist.GenerateETag(data)
	if err != nil {
		return etag, fmt.Errorf("failed to generate etag: %w", err)
	}

	if newEtag != etag {
		etag, err = persist.Store.Set(data.MainIdentity, data)
		if err != nil {
			return etag, fmt.Errorf("failed to update store: %w", err)
		}
	}

	session.Values[handlers.PreviousEtagUsed] = etag
	if authenticatedUsers, ok := session.Values[handlers.AllAuthenticatedCharacters].([]int64); ok {
		session.Values[handlers.PreviousUserCount] = len(authenticatedUsers)
	}

	if err := session.Save(r, w); err != nil {
		return etag, fmt.Errorf("failed to save session: %w", err)
	}

	return etag, nil
}
