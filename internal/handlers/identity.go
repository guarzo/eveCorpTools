package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gorilla/sessions"

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/service"
	"github.com/guarzo/zkillanalytics/internal/utils"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

// ValidUser now checks if a character is explicitly listed in the configuration or is a trusted character.
func ValidUser(character model.CharacterData) bool {
	return slices.Contains(config.CharacterIDs, int(character.CharacterID)) ||
		IsTrustedCharacter(character.CharacterID)
}

// IsTrustedCharacter checks if the character is in the list of trusted characters, ignoring corporations.
func IsTrustedCharacter(characterID int64) bool {
	trustedCharacters, _ := persist.LoadTrustedCharacters()

	for _, char := range trustedCharacters.TrustedCharacters {
		if char.CharacterID == characterID {
			return true
		}
	}
	return false
}

func SameUserCount(session *sessions.Session, storeUsers int) bool {
	sessionValues := GetSessionValues(session)

	if sessionValues.PreviousUserCount == 0 {
		return false
	}

	if sessionValues.PreviousUserCount != storeUsers {
		return false
	}

	if authenticatedUsers, ok := session.Values[AllAuthenticatedCharacters].([]int64); ok {
		return sessionValues.PreviousUserCount == len(authenticatedUsers)
	}

	return false
}

func CheckIfCanSkip(session *sessions.Session) (model.StoreData, string, bool) {
	canSkip := true
	sessionValues := GetSessionValues(session)
	storeData, etag, ok := persist.Store.Get(sessionValues.LoggedInUser)
	if !ok {
		canSkip = false
	}

	if !SameUserCount(session, len(storeData.Identities)) {
		canSkip = false

	}
	return storeData, etag, canSkip
}

func ResetIdentitiesHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, SessionName)
		mainIdentity, ok := session.Values[LoggedInUser].(int64)
		host := utils.GetHost(r.Host)

		if !ok || mainIdentity == 0 {
			handleAuthErrorWithRedirect(w, r, "Attempt to reset identities without a main identity", "/logout")
			return
		}

		err := persist.DeleteIdentity(mainIdentity, host)
		if err != nil {
			xlog.Logf("Failed to delete identity %d: %v", mainIdentity, err)
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func SameIdentities(users []int64, identities map[int64]model.CharacterData) bool {
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

func ValidateIdentities(session *sessions.Session, esiService *service.EsiService, r *http.Request, w http.ResponseWriter) (map[int64]model.CharacterData, error) {
	sessionValues := GetSessionValues(session)
	storeData, etag, canSkip := CheckIfCanSkip(session)
	identities := storeData.Identities
	host := utils.GetHost(r.Host)
	xlog.Logf("validate identities for %s", host)

	if !canSkip {
		authenticatedUsers, ok := session.Values[AllAuthenticatedCharacters].([]int64)
		if !ok {
			xlog.Logf("Failed to retrieve authenticated users from session")
			return nil, fmt.Errorf("failed to retrieve authenticated users from session")
		}

		needIdentityPopulation := len(authenticatedUsers) == 0 || !SameIdentities(authenticatedUsers, storeData.Identities) || time.Since(time.Unix(sessionValues.LastRefreshTime, 0)) > 15*time.Minute

		if needIdentityPopulation {
			userConfig, err := persist.LoadIdentities(sessionValues.LoggedInUser, host)

			if err != nil {
				xlog.Logf("Failed to load identities: %v", err)
				return nil, fmt.Errorf("failed to load identities: %w", err)
			}

			identities, err = esiService.EsiClient.PopulateIdentities(userConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to populate identities: %w", err)
			}

			if !ValidUser(identities[sessionValues.LoggedInUser]) {
				return nil, fmt.Errorf("not a valid user - ask in discord if you think this is a mistake")
			}

			if err = persist.SaveIdentities(sessionValues.LoggedInUser, userConfig, host); err != nil {
				return nil, fmt.Errorf("failed to save identities: %w", err)
			}

			session.Values[AllAuthenticatedCharacters] = GetAuthenticatedCharacterIDs(identities)
			session.Values[LastRefreshTime] = time.Now().Unix()
		}

		_, err := UpdateAndStoreSession(storeData, etag, session, r, w)
		if err != nil {
			xlog.Logf("Failed to update store and session: %v", err)
			return nil, errors.New("failed to update and store session")
		}
	}

	return identities, nil
}
