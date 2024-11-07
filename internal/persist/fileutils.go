package persist

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJSONFromFile reads JSON data from a file and populates the provided structure.
func ReadJSONFromFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}
	return nil
}

// WriteJSONToFile writes the given structure as JSON to a file.
func WriteJSONToFile(filename string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %v", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}
	return nil
}
