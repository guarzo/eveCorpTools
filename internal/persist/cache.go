// internal/persist/in_memory_cache.go

package persist

import (
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
	value      []byte
	expiration time.Time
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
	if !exists || time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set stores data in the cache with the associated key and expiration duration.
func (c *Cache) Set(key string, value []byte, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
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
			if now.After(item.expiration) {
				delete(c.data, key)
				c.Logger.Infof("Cache item expired and removed: %s", key)
			}
		}
		c.mutex.Unlock()
	}
}
