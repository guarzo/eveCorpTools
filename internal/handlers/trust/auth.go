package trust

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/guarzo/zkillanalytics/internal/handlers"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/service"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

func LoginHandler(esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := fmt.Sprintf("main-%d", time.Now().UnixNano())
		url := esiService.EsiClient.GetAuthURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func AuthCharacterHandler(esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := fmt.Sprintf("character-%d", time.Now().UnixNano())
		url := esiService.EsiClient.GetAuthURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// CallbackHandler handles the OAuth callback
func CallbackHandler(s *handlers.SessionService, esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		xlog.Logf("Received OAuth callback with code: %s, state: %s", code, state)

		token, err := esiService.EsiClient.ExchangeCode(code)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to exchange token: %v", err), "/")
			return
		}

		// Log only non-sensitive token information
		xlog.Logf("Exchanged code for token: TokenType=%s, Expiry=%s", token.TokenType, token.Expiry)

		// Get user information
		user, err := esiService.EsiClient.GetUserInfo(token)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to get user info: %v", err), "/")
			return
		}

		xlog.Logf("Authenticated user: CharacterID %d, Name %s", user.CharacterID, user.CharacterName)

		session, err := s.Get(r, handlers.SessionName)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to get session: %v", err), "/")
			return
		}

		if strings.HasPrefix(state, "main") {
			// Handle main identity
			session.Values[handlers.LoggedInUser] = user.CharacterID

			// Initialize all_authenticated_characters with main CharacterID
			session.Values[handlers.AllAuthenticatedCharacters] = []int64{user.CharacterID}

			xlog.Logf("Set mainIdentity to CharacterID: %d and initialized allAuthenticatedCharacters", user.CharacterID)
		} else if strings.HasPrefix(state, "character") {
			// Handle additional character
			existingChars, ok := session.Values[handlers.AllAuthenticatedCharacters].([]int64)
			if !ok {
				session.Values[handlers.AllAuthenticatedCharacters] = []int64{user.CharacterID}
				xlog.Logf("Initialized allAuthenticatedCharacters with CharacterID: %d", user.CharacterID)
			} else {
				if !contains(existingChars, user.CharacterID) {
					session.Values[handlers.AllAuthenticatedCharacters] = append(existingChars, user.CharacterID)
					xlog.Logf("Added CharacterID %d to allAuthenticatedCharacters", user.CharacterID)
				} else {
					xlog.Logf("CharacterID %d already exists in allAuthenticatedCharacters", user.CharacterID)
				}
			}

			// Log the authenticated character for additional authentication
			xlog.Logf("Authenticated character during 'character' state: %d, Name: %s", user.CharacterID, user.CharacterName)
		} else {
			// Unknown state
			handleErrorWithRedirect(w, r, "Invalid state parameter", "/")
			return
		}

		// Log session values for debugging
		xlog.Logf("Session Values after setting: logged_in_user=%v, all_authenticated_characters=%v", session.Values[handlers.LoggedInUser], session.Values[handlers.AllAuthenticatedCharacters])

		mainIdentity, ok := session.Values[handlers.LoggedInUser].(int64)
		if !ok || mainIdentity == 0 {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Main identity not found, current session: %v", session.Values), "/logout")
			return
		}

		xlog.Logf("MainIdentity: %d", mainIdentity)

		// Update identities with the new token
		err = persist.UpdateIdentities(mainIdentity, func(userConfig *model.Identities) error {
			xlog.Logf("Updating token for CharacterID: %d", user.CharacterID)
			userConfig.Tokens[fmt.Sprintf("%d", user.CharacterID)] = *token
			return nil
		})

		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to update identities: %v", err), "/")
			return
		}

		// Retrieve all authenticated characters from the session
		allChars, ok := session.Values[handlers.AllAuthenticatedCharacters].([]int64)
		if !ok || len(allChars) == 0 {
			xlog.Logf("Logged in characters: <nil>")
		} else {
			xlog.Logf("Logged in characters: %v", allChars)
		}

		// Save the session
		if err := session.Save(r, w); err != nil {
			xlog.Logf("Failed to save session: %v", err)
			handleErrorWithRedirect(w, r, "Failed to save session", "/")
			return
		}
		xlog.Log("Session saved successfully.")

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func LogoutHandler(s *handlers.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, handlers.SessionName)
		clearSession(s, w, r)
		session.Save(r, w)

		// Capture the 'error' query parameter if present
		errorMessage := r.URL.Query().Get("error")
		if errorMessage != "" {
			// URL-encode the error message
			encodedError := url.QueryEscape(errorMessage)
			// Append to redirect URL
			redirectURL := fmt.Sprintf("/?error=%s", encodedError)
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}

		// Otherwise, redirect normally
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func ResetIdentitiesHandler(s *handlers.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, handlers.SessionName)
		mainIdentity, ok := session.Values[handlers.LoggedInUser].(int64)

		if !ok || mainIdentity == 0 {
			handleErrorWithRedirect(w, r, "Attempt to reset identities without a main identity", "/logout")
			return
		}

		err := persist.DeleteIdentity(mainIdentity)
		if err != nil {
			xlog.Logf("Failed to delete identity %d: %v", mainIdentity, err)
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

// Helper function to check if a slice contains a specific int64 value
func contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
