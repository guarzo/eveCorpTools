package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Cache is a simple thread-safe in-memory cache.
type Cache struct {
	data   map[string]cacheItem
	mutex  sync.RWMutex
	Logger *logrus.Logger
}

type cacheItem struct {
	Value      []byte    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

// NewInMemoryCache initializes a new Cache with a provided Logger.
func NewInMemoryCache(logger *logrus.Logger) *Cache {
	cache := &Cache{
		data:   make(map[string]cacheItem),
		Logger: logger,
	}

	// Start a goroutine to clean up expired items periodically
	go cache.cleanupExpiredItems()

	return cache
}

// Get retrieves data from the cache based on the key.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}

// Set stores data in the cache with the associated key and expiration duration.
func (c *Cache) Set(key string, value []byte, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	}

	return nil
}

// cleanupExpiredItems periodically removes expired items from the cache.
func (c *Cache) cleanupExpiredItems() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		c.mutex.Lock()
		for key, item := range c.data {
			if now.After(item.Expiration) {
				delete(c.data, key)
				c.Logger.Infof("Cache item expired and removed: %s", key)
			}
		}
		c.mutex.Unlock()
	}
}

// SaveToFile saves the entire cache to a file in JSON format.
func (c *Cache) SaveToFile(filename string) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serialize the cache data to JSON and write to file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(c.data); err != nil {
		return err
	}

	c.Logger.Infof("Cache successfully saved to file: %s", filename)
	return nil
}

// LoadFromFile loads the cache from a JSON file.
func (c *Cache) LoadFromFile(filename string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			c.Logger.Warnf("Cache file does not exist: %s", filename)
			return nil // Not an error if the cache file doesn't exist
		}
		return err
	}
	defer file.Close()

	// Decode JSON data into the cache
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c.data); err != nil {
		return err
	}

	c.Logger.Infof("Cache successfully loaded from file: %s", filename)
	return nil
}

func GenerateCacheDataFileName() string {
	return fmt.Sprintf("%s/cache.json", GenerateRelativeDirectoryPath(dataDirectory))
}
