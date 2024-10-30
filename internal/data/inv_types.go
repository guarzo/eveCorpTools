package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
)

type InvTypeService struct {
	invTypes   []model.InvType
	invTypeMap map[int]string
	logger     *logrus.Logger
}

// NewInvTypeService initializes both the slice and the map.
func NewInvTypeService(logger *logrus.Logger) *InvTypeService {
	return &InvTypeService{
		invTypes:   []model.InvType{},
		invTypeMap: make(map[int]string),
		logger:     logger,
	}
}

func (iv *InvTypeService) LoadInvTypes() error {
	filePath := filepath.Join(persist.GenerateRelativeDirectoryPath("static"), "types.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open types.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // set the delimiter
	reader.TrimLeadingSpace = true

	lines, err := reader.ReadAll()
	if err != nil {
		iv.logger.Fatal(err)
	}

	for i, line := range lines {
		// Skip the header line
		if i == 0 {
			continue
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return fmt.Errorf("invalid lines in types.csv: %w", err)
		}
		iv.invTypes = append(iv.invTypes, model.InvType{
			ID:   id,
			Name: line[1],
		})
		iv.invTypeMap[id] = line[1]
	}

	iv.logger.Infof("size of inv types: %v", len(iv.invTypeMap))
	return nil
}

func (iv *InvTypeService) QueryInvType(id int) string {
	for _, invType := range iv.invTypes {
		if invType.ID == id {
			return invType.Name
		}
	}
	return "Unknown"
}
