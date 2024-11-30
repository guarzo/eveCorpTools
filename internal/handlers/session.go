package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

const (
	LastRefreshTime            = "last_refresh"
	AllAuthenticatedCharacters = "authenticated_characters"
	LoggedInUser               = "logged_in_user"
	SessionName                = "session"
	PreviousUserCount          = "previous_user_count"
	PreviousInputSubmitted     = "previous_input_submitted"
	PreviousEtagUsed           = "previous_etag_used"
)

type SessionValues struct {
	LastRefreshTime        int64
	LoggedInUser           int64
	PreviousUserCount      int
	PreviousInputSubmitted string
	PreviousEtagUsed       string
}

type SessionService struct {
	store *sessions.CookieStore
}

func GetSessionValues(session *sessions.Session) SessionValues {
	s := SessionValues{}

	if userID, ok := session.Values[LoggedInUser].(int64); ok {
		s.LoggedInUser = userID
	} else {
		xlog.Log("logged_in_user not found or not an int64 in session")
		s.LoggedInUser = 0
	}

	if val, ok := session.Values[PreviousUserCount].(int); ok {
		s.PreviousUserCount = val
	}

	if val, ok := session.Values[PreviousInputSubmitted].(string); ok {
		s.PreviousInputSubmitted = val
	}

	if val, ok := session.Values[PreviousEtagUsed].(string); ok {
		s.PreviousEtagUsed = val
	}

	if val, ok := session.Values[LastRefreshTime].(int64); ok {
		s.LastRefreshTime = val
	}

	return s
}

func NewSessionService(secret string) *SessionService {
	return &SessionService{
		store: sessions.NewCookieStore([]byte(secret)),
	}
}

func (s *SessionService) Get(r *http.Request, name string) (*sessions.Session, error) {
	return s.store.Get(r, name)
}

func ClearSession(s *SessionService, w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := s.Get(r, SessionName)
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

func UpdateAndStoreSession(data model.StoreData, etag string, session *sessions.Session, r *http.Request, w http.ResponseWriter) (string, error) {
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

	session.Values[PreviousEtagUsed] = etag
	if authenticatedUsers, ok := session.Values[AllAuthenticatedCharacters].([]int64); ok {
		session.Values[PreviousUserCount] = len(authenticatedUsers)
	}

	if err := session.Save(r, w); err != nil {
		return etag, fmt.Errorf("failed to save session: %w", err)
	}

	return etag, nil
}
