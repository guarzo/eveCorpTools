package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/service"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

var (
	Tmpl = template.Must(template.ParseFiles(
		filepath.Join("static", "tmpl", "landing.tmpl"),
	))
)

func AuthMiddleware(sessionStore *SessionService, esiService *service.EsiService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// List of public routes that don't require authentication
			publicRoutes := map[string]bool{
				"/static":   true,
				"/landing":  true,
				"/login":    true,
				"/logout":   true,
				"/callback": true,
			}

			log.Printf("Incoming request path: %s", r.URL.Path)

			// Check if the path starts with one of the public routes
			for publicRoute := range publicRoutes {
				if strings.HasPrefix(r.URL.Path, publicRoute) {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Println("Proceeding to authentication check")

			session, err := sessionStore.Get(r, SessionName)
			if err != nil {
				log.Printf("Error getting session: %v", err)
				// Clear any invalid session cookies
				http.SetCookie(w, &http.Cookie{
					Name:   SessionName,
					MaxAge: -1, // Expire the session cookie
					Path:   "/",
				})
				handleAuthErrorWithRedirect(w, r, err.Error(), "/landing")
				return
			}

			sessionValues := GetSessionValues(session)
			log.Printf("Session values: %v", sessionValues)

			// Check if logged_in_user is present
			loggedInUser := sessionValues.LoggedInUser
			if loggedInUser == 0 {
				log.Println("User not logged in, redirecting to /landing")
				http.Redirect(w, r, "/landing", http.StatusFound)
				return
			}

			// Ensure token exists for the logged-in user
			_, err = persist.GetMainIdentityToken(loggedInUser)
			if err != nil {
				// If token is missing, redirect to the landing page
				handleAuthErrorWithRedirect(w, r, err.Error(), "/landing")
				return
			}

			_, err = ValidateIdentities(session, esiService, r, w)
			if err != nil {
				xlog.Logf("Failed to validate identities")
				handleAuthErrorWithRedirect(w, r, err.Error(), "/landing")
				return
			}

			// If user is authenticated, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func GetAuthenticatedCharacterIDs(identities map[int64]model.CharacterData) []int64 {
	authenticatedCharacters := make([]int64, 0, len(identities))
	for id := range identities {
		authenticatedCharacters = append(authenticatedCharacters, id)
	}
	return authenticatedCharacters
}

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
func CallbackHandler(s *SessionService, esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		xlog.Logf("Received OAuth callback with code: %s, state: %s", code, state)

		token, err := esiService.EsiClient.ExchangeCode(code)
		if err != nil {
			handleAuthErrorWithRedirect(w, r, fmt.Sprintf("Failed to exchange token: %v", err), "/")
			return
		}

		// Log only non-sensitive token information
		xlog.Logf("Exchanged code for token: TokenType=%s, Expiry=%s", token.TokenType, token.Expiry)

		// Get user information
		user, err := esiService.EsiClient.GetUserInfo(token)
		if err != nil {
			handleAuthErrorWithRedirect(w, r, fmt.Sprintf("Failed to get user info: %v", err), "/")
			return
		}

		xlog.Logf("Authenticated user: CharacterID %d, Name %s", user.CharacterID, user.CharacterName)

		session, err := s.Get(r, SessionName)
		if err != nil {
			handleAuthErrorWithRedirect(w, r, fmt.Sprintf("Failed to get session: %v", err), "/")
			return
		}

		if strings.HasPrefix(state, "main") {
			// Handle main identity
			session.Values[LoggedInUser] = user.CharacterID

			// Initialize all_authenticated_characters with main CharacterID
			session.Values[AllAuthenticatedCharacters] = []int64{user.CharacterID}

			xlog.Logf("Set mainIdentity to CharacterID: %d and initialized allAuthenticatedCharacters", user.CharacterID)
		} else if strings.HasPrefix(state, "character") {
			// Handle additional character
			existingChars, ok := session.Values[AllAuthenticatedCharacters].([]int64)
			if !ok {
				session.Values[AllAuthenticatedCharacters] = []int64{user.CharacterID}
				xlog.Logf("Initialized allAuthenticatedCharacters with CharacterID: %d", user.CharacterID)
			} else {
				if !contains(existingChars, user.CharacterID) {
					session.Values[AllAuthenticatedCharacters] = append(existingChars, user.CharacterID)
					xlog.Logf("Added CharacterID %d to allAuthenticatedCharacters", user.CharacterID)
				} else {
					xlog.Logf("CharacterID %d already exists in allAuthenticatedCharacters", user.CharacterID)
				}
			}

			// Log the authenticated character for additional authentication
			xlog.Logf("Authenticated character during 'character' state: %d, Name: %s", user.CharacterID, user.CharacterName)
		} else {
			// Unknown state
			handleAuthErrorWithRedirect(w, r, "Invalid state parameter", "/")
			return
		}

		// Log session values for debugging
		xlog.Logf("Session Values after setting: logged_in_user=%v, all_authenticated_characters=%v", session.Values[LoggedInUser], session.Values[AllAuthenticatedCharacters])

		mainIdentity, ok := session.Values[LoggedInUser].(int64)
		if !ok || mainIdentity == 0 {
			handleAuthErrorWithRedirect(w, r, fmt.Sprintf("Main identity not found, current session: %v", session.Values), "/logout")
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
			handleAuthErrorWithRedirect(w, r, fmt.Sprintf("Failed to update identities: %v", err), "/")
			return
		}

		// Retrieve all authenticated characters from the session
		allChars, ok := session.Values[AllAuthenticatedCharacters].([]int64)
		if !ok || len(allChars) == 0 {
			xlog.Logf("Logged in characters: <nil>")
		} else {
			xlog.Logf("Logged in characters: %v", allChars)
		}

		// Save the session
		if err := session.Save(r, w); err != nil {
			xlog.Logf("Failed to save session: %v", err)
			handleAuthErrorWithRedirect(w, r, "Failed to save session", "/")
			return
		}
		xlog.Log("Session saved successfully.")

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func LogoutHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, SessionName)
		ClearSession(s, w, r)
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

// Helper function to check if a slice contains a specific int64 value
func contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
