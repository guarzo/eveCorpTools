package model

import (
	"time"
)

type Ids struct {
	AllianceIDs    []int `json:"alliance_ids"`
	CharacterIDs   []int `json:"character_ids"`
	CorporationIDs []int `json:"corporation_ids"`
}

type ZKB struct {
	LocationID     int64   `json:"locationID"`
	Hash           string  `json:"hash"`
	FittedValue    float64 `json:"fittedValue"`
	DroppedValue   float64 `json:"droppedValue"`
	DestroyedValue float64 `json:"destroyedValue"`
	TotalValue     float64 `json:"totalValue"`
	Points         int     `json:"points"`
	NPC            bool    `json:"npc"`
	Solo           bool    `json:"solo"`
	Awox           bool    `json:"awox"`
}

type KillMail struct {
	KillMailID int64 `json:"killmail_id"`
	ZKB        ZKB   `json:"zkb"`
}

type Attacker struct {
	AllianceID     int     `json:"alliance_id"`
	CharacterID    int     `json:"character_id"`
	CorporationID  int     `json:"corporation_id"`
	DamageDone     int     `json:"damage_done"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float64 `json:"security_status"`
	ShipTypeID     int     `json:"ship_type_id"`
	WeaponTypeID   int     `json:"weapon_type_id"`
}

type Victim struct {
	CharacterID   int           `json:"character_id"`
	CorporationID int           `json:"corporation_id"`
	DamageTaken   int           `json:"damage_taken"`
	Items         []interface{} `json:"items"`
	Position      struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"position"`
	ShipTypeID int `json:"ship_type_id"`
}

type EsiKillMail struct {
	KillMailID    int        `json:"killmail_id"`
	KillMailTime  time.Time  `json:"killmail_time"`
	SolarSystemID int        `json:"solar_system_id"`
	Victim        Victim     `json:"victim"`
	Attackers     []Attacker `json:"attackers"`
}

// Character contains detailed information about an EVE Online character
type Character struct {
	Birthday       time.Time `json:"birthday"`
	BloodlineID    int       `json:"bloodline_id"`
	CorporationID  int       `json:"corporation_id"`
	Description    string    `json:"description"`
	Gender         string    `json:"gender"`
	Name           string    `json:"name"`
	RaceID         int       `json:"race_id"`
	SecurityStatus float64   `json:"security_status"`
}

// Alliance contains detailed information about an EVE Online alliance
type Alliance struct {
	CreatorCorporationID  int       `json:"creator_corporation_id"`
	CreatorID             int       `json:"creator_id"`
	DateFounded           time.Time `json:"date_founded"`
	ExecutorCorporationID int       `json:"executor_corporation_id"`
	Name                  string    `json:"name"`
	Ticker                string    `json:"ticker"`
}

// Corporation represents detailed information about an EVE Online corporation
type Corporation struct {
	AllianceID    int       `json:"alliance_id"`
	CeoID         int       `json:"ceo_id"`
	CreatorID     int       `json:"creator_id"`
	DateFounded   time.Time `json:"date_founded"`
	Description   string    `json:"description"`
	HomeStationID int       `json:"home_station_id"`
	MemberCount   int       `json:"member_count"`
	Name          string    `json:"name"`
	Shares        int       `json:"shares"`
	TaxRate       float64   `json:"tax_rate"`
	Ticker        string    `json:"ticker"`
	URL           string    `json:"url"`
}

// DetailedKillMail contains both the basic and detailed killmail data
type DetailedKillMail struct {
	KillMail
	EsiKillMail
}

type KillMailData struct {
	KillMails []DetailedKillMail
}

type ESIData struct {
	AllianceInfos    map[int]Alliance
	CharacterInfos   map[int]Character
	CorporationInfos map[int]Corporation
}

type ChartData struct {
	KillMails []DetailedKillMail
	ESIData
}

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

func (c Corporation) GetName() string { return c.Name }
func (a Alliance) GetName() string    { return a.Name }
func (c Character) GetName() string   { return c.Name }

type Namer interface {
	GetName() string
}

type LootSplit struct {
	TotalBuyPrice string            `json:"totalBuyPrice"`
	SplitDetails  map[string]string `json:"splitDetails"`
	BattleReport  string            `json:"battleReport"`
	Date          string            `json:"date"`
}
