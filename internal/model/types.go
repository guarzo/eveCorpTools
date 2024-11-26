package model

import (
	"encoding/json"
	"strconv"
)

type BattleReport struct {
	BattleReportLink string `json:"battleReportLink"`
	Description      string `json:"description"`
	VideoLink        string `json:"videoLink"`
}

// InvType represents a type ID and Name for an EVE Online item
type InvType struct {
	ID   int
	Name string
}

type LootSplit struct {
	TotalBuyPrice string            `json:"totalBuyPrice"`
	SplitDetails  map[string]Amount `json:"splitDetails"` // Custom type to handle mixed input
	BattleReport  string            `json:"battleReport"`
	Date          string            `json:"date"`
	ID            int               `json:"id"`
}

type Amount float64

func (a *Amount) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as float64
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*a = Amount(num)
		return nil
	}

	// If that fails, try to unmarshal as a string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Convert the string to a float64
	parsed, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*a = Amount(parsed)
	return nil
}
