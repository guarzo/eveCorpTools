package handlers

import (
	"github.com/gorilla/sessions"
	"github.com/guarzo/zkillanalytics/internal/xlog"
	"net/http"
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
