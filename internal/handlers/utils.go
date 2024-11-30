package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/utils"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int, logger *logrus.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Errorf("Failed to encode JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal Server Error"})
	}
}

func WriteJSONError(w http.ResponseWriter, message string, identifier string, statusCode int, logger *logrus.Logger) {
	// Include the identifier in the error message if provided
	fullMessage := message
	if identifier != "" {
		fullMessage = fmt.Sprintf("%s: %s", message, identifier)
	}

	// Log the error for tracing
	logger.Errorf("Error: %s, StatusCode: %d, Identifier: %s", fullMessage, statusCode, identifier)

	// Use WriteJSONResponse to format and send the error response
	WriteJSONResponse(w, ErrorResponse{Error: fullMessage}, statusCode, logger)
}

func GetSessionIdentity(s *SessionService, r *http.Request, logger *logrus.Logger) (int64, oauth2.Token, error) {
	session, err := s.Get(r, SessionName)
	host := utils.GetHost(r.Host)
	if err != nil {
		logger.Errorf("Session retrieval error: %v", err)
		return 0, oauth2.Token{}, fmt.Errorf("failed to retrieve session")
	}

	mainIdentity, ok := session.Values[LoggedInUser].(int64)
	if !ok || mainIdentity == 0 {
		logger.Warn("Main identity not found in session")
		return 0, oauth2.Token{}, fmt.Errorf("main identity not found")
	}

	token, err := persist.GetMainIdentityToken(mainIdentity, host)
	if err != nil {
		logger.Errorf("Error retrieving token for main identity: %v", err)
		return 0, oauth2.Token{}, fmt.Errorf("failed to retrieve token")
	}

	return mainIdentity, token, nil
}

// ErrorResponse represents a JSON-formatted error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a JSON-formatted success message.
type SuccessResponse struct {
	Message string `json:"message"`
}

// HandleErrorWithRedirect redirects to the given URL with an error message as a query parameter
func HandleErrorWithRedirect(w http.ResponseWriter, r *http.Request, errorMessage, redirectURL string) {
	// URL-encode the error message to ensure it's safe for URLs
	encodedMessage := url.QueryEscape(errorMessage)

	// Construct the new URL with the error query parameter
	newURL := fmt.Sprintf("%s?error=%s", redirectURL, encodedMessage)
	xlog.Logf(newURL)

	// Redirect the user with the updated URL
	http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
}

// handleAuthErrorWithRedirect redirects to the given URL with an error message as a query parameter
func handleAuthErrorWithRedirect(w http.ResponseWriter, r *http.Request, errorMessage, redirectURL string) {

	host := utils.GetHost(r.Host)

	title := fmt.Sprintf("Zoo Auth-%s", host) // default

	// URL-encode the error message to ensure it's safe for URLs
	encodedMessage := url.QueryEscape(errorMessage)
	encodedTitle := url.QueryEscape(title)
	// Construct the new URL with the error query parameter
	newURL := fmt.Sprintf("%s?error=%s&title=%s", redirectURL, encodedMessage, encodedTitle)
	xlog.Logf(newURL)

	// Redirect the user with the updated URL
	http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
}
