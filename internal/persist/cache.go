package persist

import (
	"fmt"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

const cacheDirectory = "data/tps/cache"
const defaultExpiration = 25 * time.Hour
const cleanupInterval = 26 * time.Hour

// Cache struct to manage in-memory caching with optional persistence.
type Cache struct {
	cache  *cache.Cache
	Logger *logrus.Logger
}

func NewCache(logger *logrus.Logger) *Cache {
	return &Cache{
		cache:  cache.New(defaultExpiration, cleanupInterval),
		Logger: logger,
	}
}

// Get retrieves a value from the cache by key.
func (c *Cache) Get(key string) ([]byte, bool) {
	value, found := c.cache.Get(key)
	if !found {
		return nil, false
	}
	byteSlice, ok := value.([]byte)
	return byteSlice, ok
}

// Set stores data in the cache with the associated key and expiration duration.
func (c *Cache) Set(key string, value []byte, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

// SaveToFile saves the entire cache to a file in JSON format.
func (c *Cache) SaveToFile(filename string) error {
	items := c.cache.Items()

	// Convert to map[string]cacheItem for serialization
	serializable := make(map[string]cacheItem, len(items))
	for k, v := range items {
		byteSlice, ok := v.Object.([]byte)
		if !ok {
			c.Logger.Warnf("Skipping key %s as its value is not []byte", k)
			continue
		}
		serializable[k] = cacheItem{
			Value:      byteSlice,
			Expiration: time.Unix(0, v.Expiration),
		}
	}
	return WriteJSONToFile(filename, serializable)
}

// LoadFromFile loads the cache from a JSON file.
func (c *Cache) LoadFromFile(filename string) error {
	var serializable map[string]cacheItem
	if err := ReadJSONFromFile(filename, &serializable); err != nil {
		if os.IsNotExist(err) {
			c.Logger.Warnf("Cache file does not exist: %s", filename)
			return nil
		}
		return err
	}

	// Set each item in the cache with its expiration check
	for k, item := range serializable {
		ttl := time.Until(item.Expiration)
		if ttl > 0 {
			c.cache.Set(k, item.Value, ttl)
		} else {
			c.Logger.Infof("Skipping expired cache item: %s", k)
		}
	}
	c.Logger.Infof("Cache successfully loaded from file: %s", filename)
	return nil
}

// GenerateCacheDataFileName generates the filename for storing cache data.
func GenerateCacheDataFileName() string {
	return fmt.Sprintf("%s/cache.json", GenerateRelativeDirectoryPath(cacheDirectory))
}

// cacheItem represents an item stored in the cache with its expiration.
type cacheItem struct {
	Value      []byte
	Expiration time.Time
}
