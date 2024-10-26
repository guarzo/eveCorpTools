package esi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gambtho/zkillanalytics/internal/api"
	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/model"
)

// FetchCharacterInfo retrieves character information using the character ID
func FetchCharacterInfo(client *http.Client, characterID int) (*model.Character, error) {
	url := fmt.Sprintf("%s/characters/%d/?datasource=tranquility", config.BaseEsiURL, characterID)

	result, err := api.RetryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch character data: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var character model.Character
		if err := json.Unmarshal(body, &character); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}
		return &character, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Character), nil

}

// GetEsiKillMail retrieves detailed killmail information using the killmail ID and hash
func GetEsiKillMail(baseURL string, client *http.Client, killMailID int64, hash string) (*model.EsiKillMail, error) {
	url := fmt.Sprintf("%s/%d/%s/?datasource=tranquility", baseURL, killMailID, hash)

	result, err := api.RetryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch full killmail data: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var fullKillMail model.EsiKillMail
		if err := json.Unmarshal(body, &fullKillMail); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}

		return &fullKillMail, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.EsiKillMail), nil
}

// GetAllianceInfo retrieves alliance information using the alliance ID
func GetAllianceInfo(client *http.Client, allianceID int) (*model.Alliance, error) {
	url := fmt.Sprintf("%s/alliances/%d/?datasource=tranquility", config.BaseEsiURL, allianceID)

	result, err := api.RetryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch alliance data: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var alliance model.Alliance
		if err := json.Unmarshal(body, &alliance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}

		return &alliance, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Alliance), nil
}

func GetCorporationInfo(client *http.Client, corporationID int) (*model.Corporation, error) {
	url := fmt.Sprintf("%s/corporations/%d/?datasource=tranquility", config.BaseEsiURL, corporationID)

	result, err := api.RetryWithExponentialBackoff(func() (interface{}, error) {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch corporation data: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %s", err)
		}

		var corporation model.Corporation
		if err := json.Unmarshal(body, &corporation); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}

		return &corporation, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Corporation), nil
}
