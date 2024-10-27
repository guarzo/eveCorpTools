package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gambtho/zkillanalytics/internal/api/esi"
	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

// EsiService encapsulates business logic for ESI-related operations.
type EsiService struct {
	EsiClient *esi.EsiClient
	Cache     *persist.Cache
	Logger    *logrus.Logger
}

// NewEsiService initializes and returns a new instance of EsiService.
func NewEsiService(esiClient *esi.EsiClient, cache *persist.Cache, logger *logrus.Logger) *EsiService {
	return &EsiService{
		EsiClient: esiClient,
		Cache:     cache,
		Logger:    logger,
	}
}

// GetKillMail retrieves killmail details using the EsiClient.
func (es *EsiService) GetKillMail(ctx context.Context, killMailID int, hash string) (*model.EsiKillMail, error) {
	es.Logger.Infof("Fetching killmail ID: %d, Hash: %s", killMailID, hash)

	// Fetch the killmail details using EsiClient.
	killMail, err := es.EsiClient.GetEsiKillMail(ctx, killMailID, hash)
	if err != nil {
		es.Logger.Errorf("Error fetching killmail: %v", err)
		return nil, err
	}

	return killMail, nil
}

// GetCorporationInfo retrieves detailed information about a corporation.
func (es *EsiService) GetCorporationInfo(ctx context.Context, corporationID int) (*model.Corporation, error) {
	es.Logger.Infof("Fetching corporation ID: %d", corporationID)

	// Fetch corporation details using EsiClient.
	corp, err := es.EsiClient.GetCorporationInfo(ctx, corporationID)
	if err != nil {
		es.Logger.Errorf("Error fetching corporation info: %v", err)
		return nil, err
	}

	return corp, nil
}

// GetCharacterInfo retrieves detailed information about a character.
func (es *EsiService) GetCharacterInfo(ctx context.Context, characterID int) (*model.Character, error) {
	es.Logger.Infof("Fetching character ID: %d", characterID)

	// Fetch character details using EsiClient.
	char, err := es.EsiClient.GetCharacterInfo(ctx, characterID)
	if err != nil {
		es.Logger.Errorf("Error fetching character info: %v", err)
		return nil, err
	}

	return char, nil
}

// GetAllianceInfo retrieves detailed information about an alliance.
func (es *EsiService) GetAllianceInfo(ctx context.Context, allianceID int) (*model.Alliance, error) {
	es.Logger.Infof("Fetching alliance ID: %d", allianceID)

	// Fetch alliance details using EsiClient.
	alliance, err := es.EsiClient.GetAllianceInfo(ctx, allianceID)
	if err != nil {
		es.Logger.Errorf("Error fetching alliance info: %v", err)
		return nil, err
	}

	return alliance, nil
}

// AggregateEsiData aggregates additional ESI data related to a killmail.
// This method populates the provided ESIData structure with relevant information.
func (es *EsiService) AggregateEsiData(ctx context.Context, killMail *model.EsiKillMail, esiData *model.ESIData) error {
	es.Logger.Infof("Aggregating ESI data for killmail ID: %d", killMail.KillMailID)

	// Handle Corporation Information
	if killMail.Victim.CorporationID != 0 {
		if _, exists := esiData.CorporationInfos[killMail.Victim.CorporationID]; !exists {
			corpData, err := es.GetCorporationInfo(ctx, killMail.Victim.CorporationID)
			if err != nil {
				return fmt.Errorf("failed to fetch corporation details: %w", err)
			}
			esiData.CorporationInfos[killMail.Victim.CorporationID] = *corpData
		}
	}

	// Handle Character Information
	if killMail.Victim.CharacterID != 0 {
		if _, exists := esiData.CharacterInfos[killMail.Victim.CharacterID]; !exists {
			charData, err := es.GetCharacterInfo(ctx, killMail.Victim.CharacterID)
			if err != nil {
				return fmt.Errorf("failed to fetch character details: %w", err)
			}
			esiData.CharacterInfos[killMail.Victim.CharacterID] = *charData
		}
	}

	// Handle Alliance Information for Attackers
	for _, attacker := range killMail.Attackers {
		if attacker.AllianceID != 0 {
			if _, exists := esiData.AllianceInfos[attacker.AllianceID]; !exists {
				allianceData, err := es.GetAllianceInfo(ctx, attacker.AllianceID)
				if err != nil {
					return fmt.Errorf("failed to fetch alliance details: %w", err)
				}
				esiData.AllianceInfos[attacker.AllianceID] = *allianceData
			}
		}

		// Optionally handle Corporation and Character Information for Attackers
		if attacker.CorporationID != 0 && attacker.CorporationID != killMail.Victim.CorporationID {
			if _, exists := esiData.CorporationInfos[attacker.CorporationID]; !exists {
				corpData, err := es.GetCorporationInfo(ctx, attacker.CorporationID)
				if err != nil {
					return fmt.Errorf("failed to fetch corporation details for attacker: %w", err)
				}
				esiData.CorporationInfos[attacker.CorporationID] = *corpData
			}
		}

		if attacker.CharacterID != 0 {
			if _, exists := esiData.CharacterInfos[attacker.CharacterID]; !exists {
				charData, err := es.GetCharacterInfo(ctx, attacker.CharacterID)
				if err != nil {
					return fmt.Errorf("failed to fetch character details for attacker: %w", err)
				}
				esiData.CharacterInfos[attacker.CharacterID] = *charData
			}
		}
	}

	es.Logger.Infof("Successfully aggregated ESI data for killmail ID: %d", killMail.KillMailID)
	return nil
}

// RefreshEsiData refreshes character information in ChartData.
func (es *EsiService) RefreshEsiData(ctx context.Context, chartData *model.ChartData, client *http.Client) error {
	es.Logger.Info("Refreshing character information in ChartData.")

	for characterID := range chartData.ESIData.CharacterInfos {
		es.Logger.Infof("Refreshing character ID: %d", characterID)
		charData, err := es.GetCharacterInfo(ctx, characterID)
		if err != nil {
			es.Logger.Errorf("Failed to refresh character ID %d: %v", characterID, err)
			continue
		}
		chartData.ESIData.CharacterInfos[characterID] = *charData
	}

	es.Logger.Info("Character information refresh complete.")
	return nil
}
