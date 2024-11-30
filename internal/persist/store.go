package persist

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/guarzo/zkillanalytics/internal/model"
)

var Store *DataStore

// DataStore stores StoreData in memory
type DataStore struct {
	sync.RWMutex
	store map[int64]model.StoreData
	ETag  string
}

// NewHomeDataStore creates a new DataStore
func NewHomeDataStore() *DataStore {
	return &DataStore{
		store: make(map[int64]model.StoreData),
		ETag:  "",
	}
}

func (s *DataStore) Set(id int64, storeData model.StoreData) (string, error) {
	s.Lock()
	defer s.Unlock()
	s.store[id] = storeData

	etag, err := GenerateETag(storeData)
	if err != nil {
		return "", err
	}
	s.ETag = etag
	return etag, nil
}

func (s *DataStore) Get(id int64) (model.StoreData, string, bool) {
	s.RLock()
	defer s.RUnlock()
	homeData, ok := s.store[id]
	return homeData, s.ETag, ok
}

// Delete removes an identity from the store
func (s *DataStore) Delete(id int64) {
	s.Lock()
	defer s.Unlock()
	delete(s.store, id)
}

func GenerateETag(storeData model.StoreData) (string, error) {
	data, err := json.Marshal(storeData)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func init() {
	Store = NewHomeDataStore()
}
