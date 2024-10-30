package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// Cache is a wrapper around patrickmn/go-cache to maintain the same interface.
type Cache struct {
	cache  *cache.Cache
	Logger *logrus.Logger
}

// cacheItem represents the serialized form of a cache item.
type cacheItem struct {
	Value      []byte    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

// Constants for cache expiration handling
const (
	// cacheNoExpirationInt64 represents the NoExpiration value as int64 for comparison.
	cacheNoExpirationInt64 = int64(cache.NoExpiration)
)

// NewInMemoryCache initializes a new Cache with a provided Logger.
func NewInMemoryCache(logger *logrus.Logger) *Cache {
	// Initialize go-cache with no default expiration and a cleanup interval of 1 hour
	c := cache.New(cache.NoExpiration, 1*time.Hour)

	return &Cache{
		cache:  c,
		Logger: logger,
	}
}

// Get retrieves data from the cache based on the key.
func (c *Cache) Get(key string) ([]byte, bool) {
	value, found := c.cache.Get(key)
	if !found {
		return nil, false
	}

	byteSlice, ok := value.([]byte)
	if !ok {
		// Handle unexpected type
		c.Logger.Warnf("Cache value for key %s is not []byte", key)
		return nil, false
	}
	return byteSlice, true
}

// Set stores data in the cache with the associated key and expiration duration.
func (c *Cache) Set(key string, value []byte, expiration time.Duration) error {
	c.cache.Set(key, value, expiration)
	return nil
}

// SaveToFile saves the entire cache to a file in JSON format.
func (c *Cache) SaveToFile(filename string) error {
	items := c.cache.Items()

	// Convert to map[string]cacheItem for serialization
	serializable := make(map[string]cacheItem, len(items))
	for k, v := range items {
		// Ensure the stored value is []byte
		byteSlice, ok := v.Object.([]byte)
		if !ok {
			c.Logger.Warnf("Skipping key %s as its value is not []byte", k)
			continue
		}

		var expiration time.Time
		if v.Expiration == cacheNoExpirationInt64 {
			// Represent no expiration as zero value
			expiration = time.Time{}
		} else {
			// Convert Unix nanoseconds to time.Time
			expiration = time.Unix(0, v.Expiration)
		}

		serializable[k] = cacheItem{
			Value:      byteSlice,
			Expiration: expiration,
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: for pretty-printing
	if err := encoder.Encode(serializable); err != nil {
		return err
	}

	c.Logger.Infof("Cache successfully saved to file: %s", filename)
	return nil
}

// LoadFromFile loads the cache from a JSON file.
func (c *Cache) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			c.Logger.Warnf("Cache file does not exist: %s", filename)
			return nil // Not an error if the cache file doesn't exist
		}
		return err
	}
	defer file.Close()

	// Deserialize into map[string]cacheItem
	serializable := make(map[string]cacheItem)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&serializable); err != nil {
		return err
	}

	// Set each item in the cache
	for k, item := range serializable {
		var expiration time.Duration
		if !item.Expiration.IsZero() {
			ttl := time.Until(item.Expiration)
			if ttl > 0 {
				expiration = ttl
			} else {
				// Item already expired, skip
				c.Logger.Infof("Skipping expired cache item: %s", k)
				continue
			}
		} else {
			expiration = cache.NoExpiration
		}
		c.cache.Set(k, item.Value, expiration)
	}

	c.Logger.Infof("Cache successfully loaded from file: %s", filename)
	return nil
}

func GenerateCacheDataFileName() string {
	return fmt.Sprintf("%s/cache.json", GenerateRelativeDirectoryPath(dataDirectory))
}
