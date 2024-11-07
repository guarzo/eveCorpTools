package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"golang.org/x/oauth2"
	"strings"
)

// CharacterIDSearch searches for a character ID by name using a character ID context.
func (esi *EsiClient) CharacterIDSearch(characterID int64, name string, token *oauth2.Token) (int32, error) {
	return esi.IDSearch(characterID, name, "character", token)
}

// CorporationIDSearch searches for a corporation ID by name using a character ID context.
func (esi *EsiClient) CorporationIDSearch(characterID int64, name string, token *oauth2.Token) (int32, error) {
	return esi.IDSearch(characterID, name, "corporation", token)
}

// IDSearch performs the search logic for Character and Corporation IDs by category.
func (esi *EsiClient) IDSearch(characterID int64, name, category string, token *oauth2.Token) (int32, error) {
	// Build the search endpoint using the characterID context.
	baseURL := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/search/", characterID)
	params := map[string]string{
		"categories": category,
		"datasource": "tranquility",
		"language":   "en",
		"search":     name,
		"strict":     "true",
		"token":      token.AccessToken,
	}

	esi.Logger.Infof("Searching for %s ID by name: %s, character ID %s", category, name, characterID)

	// Generate cache key to prevent unnecessary API calls.
	cacheKey := esi.generateCacheKey(baseURL, params)
	if cachedData, found := esi.Cache.Get(cacheKey); found {
		var result map[string][]int32
		if err := json.Unmarshal(cachedData, &result); err == nil {
			if ids, exists := result[category]; exists && len(ids) > 0 {
				return ids[0], nil
			}
		}
	}

	esi.Logger.Infof("Searching for %s ID by name: %s, character ID %d, baseURL %s", category, name, characterID, baseURL)
	// Execute the request and handle response.
	bodyBytes, err := esi.getEsiEntityWithToken(baseURL, token, params)
	if err != nil {
		return 0, err
	}
	esi.Cache.Set(cacheKey, bodyBytes, defaultCacheExpiration)

	// Parse the JSON response.
	var result map[string][]int32
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return 0, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Retrieve the ID list based on category and verify results.
	ids, exists := result[category]
	if !exists || len(ids) == 0 {
		return 0, fmt.Errorf("no IDs returned for that name")
	}

	// Handle multiple IDs if necessary, verifying the exact match by name if needed.
	tempID := ids[0]
	if len(ids) > 1 {
		found := false
		for _, id := range ids {
			data, err := esi.GetPublicCharacterData(int64(id), token)
			if err != nil {
				continue
			}
			if strings.EqualFold(data.Name, name) {
				tempID = id
				found = true
				break
			}
		}
		if !found {
			return 0, fmt.Errorf("invalid IDs returned for that name")
		}
	}

	return tempID, nil
}

// GetPublicCharacterData fetches public character data.
func (esi *EsiClient) GetPublicCharacterData(characterID int64, token *oauth2.Token) (*model.CharacterResponse, error) {
	return esi.GetCharacterData(characterID, token)
}

// GetCharacterData retrieves detailed character data from the ESI API.
func (esi *EsiClient) GetCharacterData(characterID int64, token *oauth2.Token) (*model.CharacterResponse, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/?datasource=tranquility", characterID)
	data, err := esi.getEsiEntityWithToken(url, token)
	if err != nil {
		return nil, err
	}

	var character model.CharacterResponse
	if err := json.Unmarshal(data, &character); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &character, nil
}

// GetCharacterCorporation retrieves the corporation ID for a character.
func (esi *EsiClient) GetCharacterCorporation(characterID int64, token *oauth2.Token) (int32, error) {
	data, err := esi.GetCharacterData(characterID, token)
	if err != nil {
		return 0, err
	}

	return data.CorporationID, nil
}

// GetCharacterInfo retrieves character information from the ESI API, handling 404 errors by adding the character ID to a failed list.
func (esi *EsiClient) GetCharacterInfo(ctx context.Context, characterID int) (*model.Character, error) {
	// Check if character ID is already in the failed list
	if _, exists := esi.Failed.CharacterIDs[characterID]; exists {
		esi.Logger.Warnf("Skipping character %d as it previously failed with a 404", characterID)
		return nil, &model.NotFoundError{CharacterID: characterID}
	}

	// Define the ESI endpoint
	endpoint := fmt.Sprintf("characters/%d/", characterID)

	// Initialize the character struct to store the result
	var character model.Character

	// Fetch and populate character data using getEsiEntity
	if err := esi.getEsiEntity(ctx, endpoint, &character); err != nil {
		// Handle 404 error specifically
		if strings.Contains(err.Error(), "404") {
			esi.Logger.Warnf("Character %d not found, adding to failed list", characterID)
			esi.Failed.CharacterIDs[characterID] = true

			// Save the updated failed characters list
			if saveErr := persist.SaveFailedCharacters(esi.Failed); saveErr != nil {
				esi.Logger.Errorf("Failed to save failed character ID %d: %v", characterID, saveErr)
			}
			return nil, &model.NotFoundError{CharacterID: characterID}
		}

		// For other errors, log and return as-is
		esi.Logger.Debugf("Error fetching character ID %d: %v", characterID, err)
		return nil, err
	}

	return &character, nil
}

func (esi *EsiClient) GetCharacterPortrait(characterID int64) (string, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/portrait/?datasource=tranquility", characterID)
	cacheKey := fmt.Sprintf("portrait:%d", characterID)

	if cachedData, found := esi.Cache.Get(cacheKey); found {
		var portrait model.CharacterPortrait
		if err := json.Unmarshal(cachedData, &portrait); err == nil {
			return portrait.Px64x64, nil
		}
	}

	data, err := esi.getEsiEntityWithToken(url, nil)
	if err != nil {
		return "", err
	}

	var portrait model.CharacterPortrait
	if err := json.Unmarshal(data, &portrait); err != nil {
		return "", fmt.Errorf("failed to decode response body: %v", err)
	}

	esi.Cache.Set(cacheKey, data, defaultCacheExpiration)
	return portrait.Px64x64, nil
}

// GetCorporationInfo fetches and returns corporation details.
func (esi *EsiClient) GetCorporationInfo(ctx context.Context, corporationID int) (*model.Corporation, error) {
	var corporation model.Corporation
	err := esi.getEsiEntity(ctx, fmt.Sprintf("corporations/%d/", corporationID), &corporation)
	if err != nil {
		return nil, err
	}
	return &corporation, nil
}

// GetAllianceInfo fetches and returns alliance details.
func (esi *EsiClient) GetAllianceInfo(ctx context.Context, allianceID int) (*model.Alliance, error) {
	var alliance model.Alliance
	err := esi.getEsiEntity(ctx, fmt.Sprintf("alliances/%d/", allianceID), &alliance)
	if err != nil {
		return nil, err
	}
	return &alliance, nil
}
