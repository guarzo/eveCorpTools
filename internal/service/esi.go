package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/api/esi"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
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
	// Fetch character details using EsiClient.
	char, err := es.EsiClient.GetCharacterInfo(ctx, characterID)
	if err != nil {
		var notFoundError *model.NotFoundError
		if errors.As(err, &notFoundError) {
			// Handle the 404 case silently or with a warning, without logging as an error
			es.Logger.Debugf("Character %d not found; skipping\n", characterID)

		} else {
			es.Logger.Errorf("Error fetching character info: %v", err)
		}

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

// LoadTrackedCharacters loads all tracked characters from the killmails into ESIData.
func (es *EsiService) LoadTrackedCharacters(ctx context.Context, killMails []model.DetailedKillMail, esiData *model.ESIData) error {
	es.Logger.Info("Loading tracked characters into ESIData")

	for _, km := range killMails {
		// Add victim to ESI data if not already present
		if km.Victim.CharacterID != 0 {
			if _, exists := esiData.CharacterInfos[km.Victim.CharacterID]; !exists {
				victimData, err := es.GetCharacterInfo(ctx, km.Victim.CharacterID)
				if err != nil {
					es.Logger.Errorf("Failed to fetch victim character data for ID %d: %v", km.Victim.CharacterID, err)
					continue
				}
				esiData.CharacterInfos[km.Victim.CharacterID] = *victimData
			}
		}

		// Add attackers to ESI data if not already present
		for _, attacker := range km.Attackers {
			if attacker.CharacterID != 0 {
				if _, exists := esiData.CharacterInfos[attacker.CharacterID]; !exists {
					attackerData, err := es.GetCharacterInfo(ctx, attacker.CharacterID)
					if err != nil {
						es.Logger.Errorf("Failed to fetch attacker character data for ID %d: %v", attacker.CharacterID, err)
						continue
					}
					esiData.CharacterInfos[attacker.CharacterID] = *attackerData
				}
			}

			// Add corporation and alliance info for attacker
			if attacker.CorporationID != 0 {
				if _, exists := esiData.CorporationInfos[attacker.CorporationID]; !exists {
					corpData, err := es.GetCorporationInfo(ctx, attacker.CorporationID)
					if err != nil {
						es.Logger.Errorf("Failed to fetch corporation data for ID %d: %v", attacker.CorporationID, err)
						continue
					}
					esiData.CorporationInfos[attacker.CorporationID] = *corpData
				}
			}
			if attacker.AllianceID != 0 {
				if _, exists := esiData.AllianceInfos[attacker.AllianceID]; !exists {
					allianceData, err := es.GetAllianceInfo(ctx, attacker.AllianceID)
					if err != nil {
						es.Logger.Errorf("Failed to fetch alliance data for ID %d: %v", attacker.AllianceID, err)
						continue
					}
					esiData.AllianceInfos[attacker.AllianceID] = *allianceData
				}
			}
		}
	}

	es.Logger.Info("Finished loading tracked characters into ESIData")
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
