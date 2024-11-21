package model

import (
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

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

func (c Corporation) GetName() string { return c.Name }
func (a Alliance) GetName() string    { return a.Name }
func (c Character) GetName() string   { return c.Name }

type Namer interface {
	GetName() string
}

// FailedCharacters represents a structure to hold failed character IDs.
type FailedCharacters struct {
	CharacterIDs map[int]bool `json:"character_ids"`
}

// NotFoundError is a custom error type for representing 404 errors.
type NotFoundError struct {
	CharacterID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("character %d not found (404)", e.CharacterID)
}

type Ids struct {
	AllianceIDs    []int `json:"alliance_ids"`
	CharacterIDs   []int `json:"character_ids"`
	CorporationIDs []int `json:"corporation_ids"`
}

// TrustCharacter represents the user information
type TrustCharacter struct {
	User
	CorporationID int64  `json:"CorporationID"`
	Portrait      string `json:"Portrait"`
}

// CharacterData structure
type CharacterData struct {
	Token oauth2.Token
	TrustCharacter
}

// User represents the user information returned by the EVE SSO
type User struct {
	CharacterID   int64  `json:"CharacterID"`
	CharacterName string `json:"CharacterName"`
}

type CharacterResponse struct {
	AllianceID     int32     `json:"alliance_id,omitempty"`
	Birthday       time.Time `json:"birthday"`
	BloodlineID    int32     `json:"bloodline_id"`
	CorporationID  int32     `json:"corporation_id"`
	Description    string    `json:"description,omitempty"`
	FactionID      int32     `json:"faction_id,omitempty"`
	Gender         string    `json:"gender"`
	Name           string    `json:"name"`
	RaceID         int32     `json:"race_id"`
	SecurityStatus float64   `json:"security_status,omitempty"`
	Title          string    `json:"title,omitempty"`
}

type TrustedCharacter struct {
	CharacterID     int64     `json:"CharacterID"`
	CharacterName   string    `json:"CharacterName"`
	CorporationID   int64     `json:"CorporationID"`
	CorporationName string    `json:"CorporationName"`
	AddedBy         string    `json:"AddedBy"`
	DateAdded       time.Time `json:"DateAdded"`
	Comment         string    `json:"Comment"`
	IsOnCouch       bool      `json:"IsOnCouch"`
}

type TrustedCorporation struct {
	CorporationID   int64     `json:"CorporationID"`
	CorporationName string    `json:"CorporationName"`
	AllianceName    string    `json:"AllianceName"`
	AllianceID      int64     `json:"AllianceID"`
	DateAdded       time.Time `json:"DateAdded"`
	AddedBy         string    `json:"AddedBy"`
	Comment         string    `json:"Comment"`
	IsOnCouch       bool      `json:"IsOnCouch"`
}

type TrustedCharacters struct {
	TrustedCharacters     []TrustedCharacter   `json:"characters"`
	TrustedCorporations   []TrustedCorporation `json:"corporations"`
	UntrustedCharacters   []TrustedCharacter   `json:"untrusted_characters"`
	UntrustedCorporations []TrustedCorporation `json:"untrusted_corporations"`
}

// CharacterSearchResponse represents the array of character IDs returned from the search
type CharacterSearchResponse struct {
	CharacterIDs []int32 `json:"get_characters_character_id_search_character"`
}

type CorporationInfo struct {
	AllianceID    *int32  `json:"alliance_id,omitempty"`     // CorporationID of the alliance, if any
	CEOId         int32   `json:"ceo_id"`                    // CEO CorporationID (required)
	CreatorID     int32   `json:"creator_id"`                // Creator CorporationID (required)
	DateFounded   *string `json:"date_founded,omitempty"`    // Date the corporation was founded
	Description   *string `json:"description,omitempty"`     // CorporationID description
	FactionID     *int32  `json:"faction_id,omitempty"`      // Faction CorporationID, if any
	HomeStationID *int32  `json:"home_station_id,omitempty"` // Home station CorporationID, if any
	MemberCount   int32   `json:"member_count"`              // Number of members (required)
	Name          string  `json:"name"`                      // Full name of the corporation (required)
	Shares        *int64  `json:"shares,omitempty"`          // Number of shares, if any
	TaxRate       float64 `json:"tax_rate"`                  // Tax rate (required, float with max 1.0 and min 0.0)
	Ticker        string  `json:"ticker"`                    // Short name of the corporation (required)
	URL           *string `json:"url,omitempty"`             // CorporationID URL, if any
	WarEligible   *bool   `json:"war_eligible,omitempty"`    // War eligibility, if any
}

type CharacterPortrait struct {
	Px128x128 string `json:"px128x128"`
	Px256x256 string `json:"px256x256"`
	Px512x512 string `json:"px512x512"`
	Px64x64   string `json:"px64x64"`
}
