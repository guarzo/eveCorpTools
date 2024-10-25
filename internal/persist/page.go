package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type EntityLastPage struct {
	CorporationPages map[int]int
	AlliancePages    map[int]int
	CharacterPages   map[int]int
	LastFetchDate    time.Time
}

func NewEntityLastPage() *EntityLastPage {
	return &EntityLastPage{
		CorporationPages: make(map[int]int),
		AlliancePages:    make(map[int]int),
		CharacterPages:   make(map[int]int),
		LastFetchDate:    time.Now(),
	}
}

func (elp *EntityLastPage) GetLastPage(entityType string, entityID int) int {
	switch entityType {
	case "corporation":
		if p, ok := elp.CorporationPages[entityID]; ok {
			return p
		}
	case "alliance":
		if p, ok := elp.AlliancePages[entityID]; ok {
			return p
		}
	case "character":
		if p, ok := elp.CharacterPages[entityID]; ok {
			return p
		}
	}
	return 1 // return first page if not found
}

func (elp *EntityLastPage) UpdateLastPage(entityType string, entityID, page int) {
	switch entityType {
	case "corporation":
		elp.CorporationPages[entityID] = page
	case "alliance":
		elp.AlliancePages[entityID] = page
	case "character":
		elp.CharacterPages[entityID] = page
	}
	elp.LastFetchDate = time.Now()
}

func (elp *EntityLastPage) GetLastFetchDate() time.Time {
	return elp.LastFetchDate
}

func getPageFileName() string {
	return fmt.Sprintf("data/monthly/lastFetchPages_%s.json", time.Now().Format("2006_01"))
}

func SaveLastPage(data *EntityLastPage) error {
	fileName := getPageFileName()
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func LoadLastPage() (*EntityLastPage, error) {
	fileName := getPageFileName()
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data EntityLastPage
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return &data, err
}

func ClearLastPage() error {
	filepath := getPageFileName()
	return os.Remove(filepath)
}
