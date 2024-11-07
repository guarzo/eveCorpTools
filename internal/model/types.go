package model

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
	SplitDetails  map[string]string `json:"splitDetails"`
	BattleReport  string            `json:"battleReport"`
	Date          string            `json:"date"`
}
