package fetch

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

func FetchInvTypes() ([]model.InvType, error) {
	file, err := os.Open(persist.GenerateRelativeDirectoryPath("static") + string(os.PathSeparator) + "invTypes.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // set the delimiter
	reader.TrimLeadingSpace = true

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var invTypes []model.InvType
	for i, line := range lines {
		// Skip the header line
		if i == 0 {
			continue
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}
		invTypes = append(invTypes, model.InvType{
			ID:   id,
			Name: line[1],
		})
	}

	return invTypes, nil
}

var invTypes []model.InvType

func init() {
	var err error
	invTypes, err = FetchInvTypes()
	if err != nil {
		log.Fatal(fmt.Sprintf("Fatal error loading inventory types: %s", err))
	}
}

func QueryInvType(id int) string {
	for _, invType := range invTypes {
		if invType.ID == id {
			return invType.Name
		}
	}
	return "Unknown"
}
