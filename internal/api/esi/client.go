// internal/api/esi/esi.go

package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

// Define a default cache expiration duration.
const defaultCacheExpiration = 24 * time.Hour

// EsiClient encapsulates the HTTP client and cache for ESI API interactions.
type EsiClient struct {
	BaseURL string
	Client  *http.Client
	Cache   *persist.Cache
	Logger  *logrus.Logger
}

// NewEsiClient initializes and returns a new EsiClient.
func NewEsiClient(baseURL string, client *http.Client, cache *persist.Cache, logger *logrus.Logger) *EsiClient {
	return &EsiClient{
		BaseURL: baseURL,
		Client:  client,
		Cache:   cache,
		Logger:  logger,
	}
}

// buildRequestURL constructs the full request URL with query parameters.
func (esi *EsiClient) buildRequestURL(endpoint string, params map[string]string) (string, error) {
	base, err := url.Parse(esi.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	path, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint: %w", err)
	}

	fullURL := base.ResolveReference(path)

	// Add query parameters
	query := fullURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	fullURL.RawQuery = query.Encode()

	return fullURL.String(), nil
}

// generateCacheKey creates a unique cache key based on the endpoint and sorted parameters.
// It ensures that the key is consistent regardless of the order of parameters.
func (esi *EsiClient) generateCacheKey(endpoint string, params map[string]string) string {
	// To ensure consistency, sort the parameter keys.
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build the query string.
	queryParams := ""
	for _, key := range keys {
		queryParams += fmt.Sprintf("&%s=%s", key, params[key])
	}

	// Construct the cache key.
	cacheKey := fmt.Sprintf("esi:%s:%s", endpoint, queryParams)
	return cacheKey
}

// getEsiData is a generic method to fetch raw data for any ESI endpoint with caching.
func (esi *EsiClient) getEsiData(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {
	// Generate a unique cache key based on endpoint and sorted parameters.
	cacheKey := esi.generateCacheKey(endpoint, params)

	// Attempt to retrieve data from the cache.
	cachedData, found := esi.Cache.Get(cacheKey)
	if found {
		esi.Logger.Infof("Cache hit for key %s", cacheKey)
		return cachedData, nil
	}

	// Cache miss; construct the full URL.
	requestURL, err := esi.buildRequestURL(endpoint, params)
	if err != nil {
		esi.Logger.Errorf("Error building request URL: %v", err)
		return nil, err
	}

	esi.Logger.Infof("Fetching ESI data from URL %s", requestURL)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		esi.Logger.Errorf("Failed to create request: %v", err)
		return nil, err
	}

	// Perform the HTTP request
	resp, err := esi.Client.Do(req)
	if err != nil {
		esi.Logger.Errorf("Error making HTTP request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		esi.Logger.Errorf("Non-OK HTTP status: %s, body: %s", resp.Status, string(bodyBytes))
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		esi.Logger.Errorf("Failed to read response body: %v", err)
		return nil, err
	}

	// Cache the fetched data.
	if err := esi.Cache.Set(cacheKey, data, defaultCacheExpiration); err != nil {
		esi.Logger.Errorf("Failed to set cache for key %s: %v", cacheKey, err)
	} else {
		esi.Logger.Infof("Cached data for key %s", cacheKey)
	}

	return data, nil
}

// GetEsiKillMail fetches full killmail details with caching.
func (esi *EsiClient) GetEsiKillMail(ctx context.Context, killMailID int, hash string) (*model.EsiKillMail, error) {
	// Define the endpoint for fetching killmail details.
	endpoint := fmt.Sprintf("killmails/%d/%s/", killMailID, hash)

	// Define query parameters as required by the ESI API.
	params := map[string]string{
		"datasource": "tranquility",
		"language":   "en-us",
	}

	// Fetch raw data using the generic getEsiData method.
	data, err := esi.getEsiData(ctx, endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ESI killmail: %w", err)
	}

	// Unmarshal the data into EsiKillMail struct.
	var esiKillMail model.EsiKillMail
	if err := json.Unmarshal(data, &esiKillMail); err != nil {
		esi.Logger.Errorf("Failed to unmarshal ESI killmail data for ID %d: %v", killMailID, err)
		return nil, fmt.Errorf("failed to unmarshal ESI killmail data: %w", err)
	}

	return &esiKillMail, nil
}

// AggregateEsi aggregates ESI data into the provided EsiData structure with context support.
func (esi *EsiClient) AggregateEsi(ctx context.Context, killMail *model.EsiKillMail, esiData *model.ESIData) error {
	// Handle Corporation Information
	if killMail.Victim.CorporationID != 0 {
		if _, exists := esiData.CorporationInfos[killMail.Victim.CorporationID]; !exists {
			corpData, err := esi.GetCorporationInfo(ctx, killMail.Victim.CorporationID)
			if err != nil {
				return fmt.Errorf("failed to fetch corporation details: %w", err)
			}
			esiData.CorporationInfos[killMail.Victim.CorporationID] = *corpData
		}
	}

	// Handle Character Information
	if killMail.Victim.CharacterID != 0 {
		if _, exists := esiData.CharacterInfos[killMail.Victim.CharacterID]; !exists {
			charData, err := esi.GetCharacterInfo(ctx, killMail.Victim.CharacterID)
			if err != nil {
				return fmt.Errorf("failed to fetch character details: %w", err)
			}
			esiData.CharacterInfos[killMail.Victim.CharacterID] = *charData
		}
	}

	// Handle Alliance Information for Attackers
	for _, attacker := range killMail.Attackers {
		if attacker.AllianceID != 0 {
			if _, exists := esiData.AllianceInfos[attacker.AllianceID]; !exists {
				allianceData, err := esi.GetAllianceInfo(ctx, attacker.AllianceID)
				if err != nil {
					return fmt.Errorf("failed to fetch alliance details: %w", err)
				}
				esiData.AllianceInfos[attacker.AllianceID] = *allianceData
			}
		}

		// Optionally handle Corporation and Character Information for Attackers
		if attacker.CorporationID != 0 && attacker.CorporationID != killMail.Victim.CorporationID {
			if _, exists := esiData.CorporationInfos[attacker.CorporationID]; !exists {
				corpData, err := esi.GetCorporationInfo(ctx, attacker.CorporationID)
				if err != nil {
					return fmt.Errorf("failed to fetch corporation details for attacker: %w", err)
				}
				esiData.CorporationInfos[attacker.CorporationID] = *corpData
			}
		}

		if attacker.CharacterID != 0 {
			if _, exists := esiData.CharacterInfos[attacker.CharacterID]; !exists {
				charData, err := esi.GetCharacterInfo(ctx, attacker.CharacterID)
				if err != nil {
					return fmt.Errorf("failed to fetch character details for attacker: %w", err)
				}
				esiData.CharacterInfos[attacker.CharacterID] = *charData
			}
		}
	}

	return nil
}

// GetCorporationInfo fetches and returns corporation details.
func (esi *EsiClient) GetCorporationInfo(ctx context.Context, corporationID int) (*model.Corporation, error) {
	endpoint := fmt.Sprintf("corporations/%d/", corporationID)
	params := map[string]string{
		"datasource": "tranquility",
		"language":   "en-us",
	}

	data, err := esi.getEsiData(ctx, endpoint, params)
	if err != nil {
		return nil, err
	}

	var corporation model.Corporation
	if err := json.Unmarshal(data, &corporation); err != nil {
		esi.Logger.Errorf("Failed to unmarshal corporation data for ID %d: %v", corporationID, err)
		return nil, fmt.Errorf("failed to unmarshal corporation data: %w", err)
	}

	return &corporation, nil
}

// GetCharacterInfo fetches and returns character details.
func (esi *EsiClient) GetCharacterInfo(ctx context.Context, characterID int) (*model.Character, error) {
	endpoint := fmt.Sprintf("characters/%d/", characterID)
	params := map[string]string{
		"datasource": "tranquility",
		"language":   "en-us",
	}

	data, err := esi.getEsiData(ctx, endpoint, params)
	if err != nil {
		return nil, err
	}

	var character model.Character
	if err := json.Unmarshal(data, &character); err != nil {
		esi.Logger.Errorf("Failed to unmarshal character data for ID %d: %v", characterID, err)
		return nil, fmt.Errorf("failed to unmarshal character data: %w", err)
	}

	return &character, nil
}

// GetAllianceInfo fetches and returns alliance details.
func (esi *EsiClient) GetAllianceInfo(ctx context.Context, allianceID int) (*model.Alliance, error) {
	endpoint := fmt.Sprintf("alliances/%d/", allianceID)
	params := map[string]string{
		"datasource": "tranquility",
		"language":   "en-us",
	}

	data, err := esi.getEsiData(ctx, endpoint, params)
	if err != nil {
		return nil, err
	}

	var alliance model.Alliance
	if err := json.Unmarshal(data, &alliance); err != nil {
		esi.Logger.Errorf("Failed to unmarshal alliance data for ID %d: %v", allianceID, err)
		return nil, fmt.Errorf("failed to unmarshal alliance data: %w", err)
	}

	return &alliance, nil
}
