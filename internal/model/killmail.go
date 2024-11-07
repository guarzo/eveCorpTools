package model

import (
	"time"
)

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
