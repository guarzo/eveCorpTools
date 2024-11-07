// internal/api/zkill/zkill_client.go

package zkill

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/api"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
)

const zkillCacheExpiration = 770 * time.Hour // slightly more than 1 month

// ZkillClient handles interactions with the zKillboard API.
type ZkillClient struct {
	BaseURL string
	Client  *http.Client
	Cache   *persist.Cache
	Logger  *logrus.Logger
}

// NewZkillClient initializes and returns a new ZkillClient.
func NewZkillClient(baseURL string, client *http.Client, cache *persist.Cache, logger *logrus.Logger) *ZkillClient {
	return &ZkillClient{
		BaseURL: baseURL,
		Client:  client,
		Cache:   cache,
		Logger:  logger,
	}
}

func (zk *ZkillClient) getZkillData(ctx context.Context, cacheKey, requestURL string) ([]model.KillMail, error) {
	// Check if this is the current month to avoid caching
	currentYear, currentMonth, _ := time.Now().Date()
	isCurrentMonth := strings.Contains(cacheKey, fmt.Sprintf(":%d:%02d:", currentYear, currentMonth))

	if !isCurrentMonth {
		// Check if the data is in cache
		cachedData, found := zk.Cache.Get(cacheKey)
		if found {
			var killMails []model.KillMail
			if err := json.Unmarshal(cachedData, &killMails); err == nil {
				zk.Logger.Debugf("Cache hit for key %s", cacheKey)
				return killMails, nil
			}
			zk.Logger.Warnf("Failed to unmarshal cached data for key %s; refetching", cacheKey)
		}
	}

	// Cache miss or bypass, fetch data from the API
	data, err := api.GetPageData(ctx, zk.Client, requestURL)
	if err != nil {
		zk.Logger.Errorf("Error fetching data from URL %s: %v", requestURL, err)
		return nil, err
	}

	// Cache the fetched data for future use if it's not the current month
	if !isCurrentMonth {
		cachedBytes, err := json.Marshal(data)
		if err == nil {
			zk.Cache.Set(cacheKey, cachedBytes, zkillCacheExpiration)
		} else {
			zk.Logger.Errorf("Failed to cache data for key %s: %v", cacheKey, err)
		}
	}

	return data, nil
}

// fetchPageData is a helper method to fetch killmails from a specific API endpoint.
func (zk *ZkillClient) fetchPageData(ctx context.Context, apiType, entityType string, entityID, page, year, month int) ([]model.KillMail, error) {
	// Construct the request URL based on the apiType and generate a cache key.
	requestURL := fmt.Sprintf("%s/api/%s/%sID/%d/year/%d/month/%d/page/%d/",
		zk.BaseURL, apiType, entityType, entityID, year, month, page)
	// Add year and month to the cache key to ensure unique caching per month and year.
	cacheKey := fmt.Sprintf("zkill:%s:%sID:%d:%d:%02d:%d", apiType, entityType, entityID, year, month, page)

	zk.Logger.Debugf("Fetching %s from URL: %s", apiType, requestURL)
	return zk.getZkillData(ctx, cacheKey, requestURL)
}

// GetKillsPageData fetches killmails where entities are attackers.
func (zk *ZkillClient) GetKillsPageData(ctx context.Context, entityType string, entityID, page, year, month int) ([]model.KillMail, error) {
	return zk.fetchPageData(ctx, "kills", entityType, entityID, page, year, month)
}

// GetLossPageData fetches killmails where entities are victims (losses).
func (zk *ZkillClient) GetLossPageData(ctx context.Context, entityType string, entityID, page, year, month int) ([]model.KillMail, error) {
	return zk.fetchPageData(ctx, "losses", entityType, entityID, page, year, month)
}
