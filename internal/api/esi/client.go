package esi

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
)

// Define a default cache expiration duration.
const defaultCacheExpiration = 770 * time.Hour // slightly more than 1 month

type EsiClient struct {
	BaseURL     string
	Failed      *model.FailedCharacters
	Client      *http.Client
	Cache       *persist.Cache
	Logger      *logrus.Logger
	OAuthConfig *oauth2.Config
}

func NewEsiClient(baseURL string, failed *model.FailedCharacters, client *http.Client, cache *persist.Cache, logger *logrus.Logger) *EsiClient {
	return &EsiClient{
		BaseURL: baseURL,
		Client:  client,
		Failed:  failed,
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
func (esi *EsiClient) generateCacheKey(endpoint string, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	queryParams := ""
	for _, key := range keys {
		queryParams += fmt.Sprintf("&%s=%s", key, params[key])
	}

	return fmt.Sprintf("esi:%s:%s", endpoint, queryParams)
}
