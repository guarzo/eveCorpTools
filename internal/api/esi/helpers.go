package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"

	"github.com/guarzo/zkillanalytics/internal/api"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

// buildRequest constructs an HTTP GET request, optionally with query parameters and token.
func (esi *EsiClient) buildRequest(baseURL string, token *oauth2.Token, params ...map[string]string) (*http.Request, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Initialize parameters if not provided
	queryParams := make(map[string]string)
	if len(params) > 0 && params[0] != nil {
		queryParams = params[0]
	}

	// Add query parameters to the URL
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// Create the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set authorization header if a token is provided
	if token != nil && token.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getEsiEntityWithToken retrieves entity data from ESI, using an OAuth token if provided, and caches the response.
func (esi *EsiClient) getEsiEntityWithTokenNoCache(address string, token *oauth2.Token, params ...map[string]string) ([]byte, error) {
	// If params are not provided, initialize with an empty map to prevent out-of-range issues.
	queryParams := make(map[string]string)
	if len(params) > 0 && params[0] != nil {
		queryParams = params[0]
	}

	req, err := esi.buildRequest(address, token, queryParams)
	if err != nil {
		return nil, err
	}

	resp, err := esi.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	data, err := esi.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (esi *EsiClient) getEsiEntityWithToken(address string, token *oauth2.Token, params ...map[string]string) ([]byte, error) {
	// If params are not provided, initialize with an empty map to prevent out-of-range issues.
	queryParams := make(map[string]string)
	if len(params) > 0 && params[0] != nil {
		queryParams = params[0]
	} else {
		xlog.LogIndirect(fmt.Sprintf("No query parameters provided for cache ESI request, %s", address))
	}

	cacheKey := esi.generateCacheKey(address, queryParams)
	if cachedData, found := esi.Cache.Get(cacheKey); found {
		return cachedData, nil
	}

	// Define the retryable operation
	operation := func() (interface{}, error) {
		req, err := esi.buildRequest(address, token, queryParams)
		if err != nil {
			return nil, err
		}

		resp, err := esi.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		data, err := esi.handleResponse(resp)
		if err != nil {
			return nil, err
		}

		// Cache the response data with the default expiration.
		esi.Cache.Set(cacheKey, data, defaultCacheExpiration)
		return data, nil
	}

	// Call the retry function
	result, err := api.RetryWithExponentialBackoff(operation)
	if err != nil {
		return nil, err
	}

	return result.([]byte), nil
}

func (esi *EsiClient) getEsiEntity(ctx context.Context, endpoint string, entity interface{}) error {
	params := map[string]string{"datasource": "tranquility"}
	cacheKey := esi.generateCacheKey(endpoint, params)
	if cachedData, found := esi.Cache.Get(cacheKey); found {
		if err := json.Unmarshal(cachedData, entity); err == nil {
			return nil
		}
	}

	requestURL, err := esi.buildRequestURL(endpoint, params)
	if err != nil {
		return fmt.Errorf("failed to build request URL: %w", err)
	}

	// Define the retryable operation
	operation := func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		resp, err := esi.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		data, err := esi.handleResponse(resp)
		if err != nil {
			return nil, err
		}

		// Cache the result.
		esi.Cache.Set(cacheKey, data, defaultCacheExpiration)
		return data, nil
	}

	// Call the retry function
	result, err := api.RetryWithExponentialBackoff(operation)
	if err != nil {
		return err
	}

	// Unmarshal the response into the provided entity
	return json.Unmarshal(result.([]byte), entity)
}

// handleResponse processes the HTTP response and returns the response body or an error.
func (esi *EsiClient) handleResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-OK HTTP status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return bodyBytes, nil
}
