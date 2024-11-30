package esi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/oauth2"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

// InitializeOAuth configures OAuth2 with client credentials and callback URL.
func (esi *EsiClient) InitializeOAuth(clientID, clientSecret, callbackURL string) {
	esi.OAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Scopes: []string{
			"publicData",
			"esi-search.search_structures.v1",
			"esi-characters.write_contacts.v1",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
			TokenURL: "https://login.eveonline.com/v2/oauth/token",
		},
	}
}

// GetAuthURL returns the URL for OAuth2 authentication as a method on EsiClient.
func (esi *EsiClient) GetAuthURL(state string) string {
	return esi.OAuthConfig.AuthCodeURL(state)
}

// ExchangeCode exchanges the authorization code for an access token.
func (esi *EsiClient) ExchangeCode(code string) (*oauth2.Token, error) {
	return esi.OAuthConfig.Exchange(context.Background(), code)
}

// RefreshToken refreshes the OAuth token using the refresh token.
func (esi *EsiClient) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest(http.MethodPost, esi.OAuthConfig.Endpoint.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		esi.Logger.Errorf("Failed to create request to refresh token: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(esi.OAuthConfig.ClientID+":"+esi.OAuthConfig.ClientSecret)))

	resp, err := esi.Client.Do(req)
	if err != nil {
		esi.Logger.Errorf("Failed to make request to refresh token: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		esi.Logger.Errorf("Received non-OK status code %d: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("received non-OK status code %d", resp.StatusCode)
	}

	var token oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		esi.Logger.Errorf("Failed to decode response body: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &token, nil
}

// GetUserInfo retrieves user information.
func (esi *EsiClient) GetUserInfo(token *oauth2.Token) (*model.User, error) {
	baseURL := "https://login.eveonline.com/oauth/verify"
	data, err := esi.getEsiEntityWithTokenNoCache(baseURL, token)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &user, nil
}

// PopulateIdentities concurrently populates character data.
// PopulateIdentities concurrently populates character data.
func (esi *EsiClient) PopulateIdentities(userConfig *model.Identities) (map[int64]model.CharacterData, error) {
	characterData := make(map[int64]model.CharacterData)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for idStr, token := range userConfig.Tokens {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			xlog.Logf("Invalid CharacterID '%s': %v", idStr, err)
			continue
		}

		wg.Add(1)
		go func(id int64, token oauth2.Token) {
			defer wg.Done()

			charIdentity, err := esi.processIdentity(id, token, userConfig, &mu)
			if err != nil {
				xlog.Logf("Failed to process identity for character %d: %v", id, err)
				return
			}

			mu.Lock()
			characterData[id] = *charIdentity
			mu.Unlock()
		}(id, token)
	}

	wg.Wait()

	return characterData, nil
}

// processIdentity manages token refreshing and retrieves character data.
func (esi *EsiClient) processIdentity(id int64, token oauth2.Token, userConfig *model.Identities, mu *sync.Mutex) (*model.CharacterData, error) {
	newToken, err := esi.RefreshToken(token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token for character %d: %v", id, err)
	}
	token = *newToken

	mu.Lock()
	userConfig.Tokens[fmt.Sprintf("%d", id)] = token
	mu.Unlock()

	corp, err := esi.GetCharacterCorporation(id, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to get corp for character %d: %v", id, err)
	}

	user, err := esi.GetUserInfo(&token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	portrait, err := esi.GetCharacterPortrait(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get character portrait: %v", err)
	}

	character := model.TrustCharacter{
		User:          *user,
		CorporationID: int64(corp),
		Portrait:      portrait,
	}

	return &model.CharacterData{
		Token:          token,
		TrustCharacter: character,
	}, nil
}
