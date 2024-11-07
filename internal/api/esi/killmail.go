package esi

import (
	"context"
	"errors"
	"fmt"
	"github.com/guarzo/zkillanalytics/internal/model"
)

// GetEsiKillMail fetches full killmail details with caching.
func (esi *EsiClient) GetEsiKillMail(ctx context.Context, killMailID int, hash string) (*model.EsiKillMail, error) {
	// Define the endpoint for fetching killmail details.
	endpoint := fmt.Sprintf("killmails/%d/%s/", killMailID, hash)

	// Define the entity struct to populate.
	var esiKillMail model.EsiKillMail

	// Fetch entity data using fetchEsiEntity.
	if err := esi.getEsiEntity(ctx, endpoint, &esiKillMail); err != nil {
		return nil, fmt.Errorf("failed to fetch ESI killmail: %w", err)
	}

	return &esiKillMail, nil
}

// AggregateEsi aggregates ESI data into the provided EsiData structure with context support.
func (esi *EsiClient) AggregateEsi(ctx context.Context, killMail *model.EsiKillMail, esiData *model.ESIData) error {
	// Handle Corporation Information
	if killMail.Victim.CorporationID != 0 {
		if _, exists := esiData.CorporationInfos[killMail.Victim.CorporationID]; !exists {
			corpData, err := esi.GetCorporationInfo(ctx, killMail.Victim.CorporationID)
			if err != nil {
				return fmt.Errorf("failed to fetch corporation details: %w", err)
			}
			esiData.CorporationInfos[killMail.Victim.CorporationID] = *corpData
		}
	}

	// Handle Character Information
	if killMail.Victim.CharacterID != 0 {
		if _, exists := esiData.CharacterInfos[killMail.Victim.CharacterID]; !exists {
			charData, err := esi.GetCharacterInfo(ctx, killMail.Victim.CharacterID)
			if err != nil {
				var notFoundError *model.NotFoundError
				if errors.As(err, &notFoundError) {
					// Handle the 404 case silently or with a warning, without logging as an error
					// esi.Logger.Debugf("Character %d not found; skipping\n", killMail.Victim.CharacterID)
					return nil
				}
				return fmt.Errorf("failed to fetch character details: %w", err)
			}
			esiData.CharacterInfos[killMail.Victim.CharacterID] = *charData
		}
	}

	// Handle Alliance Information for Attackers
	for _, attacker := range killMail.Attackers {
		if attacker.AllianceID != 0 {
			if _, exists := esiData.AllianceInfos[attacker.AllianceID]; !exists {
				allianceData, err := esi.GetAllianceInfo(ctx, attacker.AllianceID)
				if err != nil {
					return fmt.Errorf("failed to fetch alliance details: %w", err)
				}
				esiData.AllianceInfos[attacker.AllianceID] = *allianceData
			}
		}

		// Optionally handle Corporation and Character Information for Attackers
		if attacker.CorporationID != 0 && attacker.CorporationID != killMail.Victim.CorporationID {
			if _, exists := esiData.CorporationInfos[attacker.CorporationID]; !exists {
				corpData, err := esi.GetCorporationInfo(ctx, attacker.CorporationID)
				if err != nil {
					return fmt.Errorf("failed to fetch corporation details for attacker: %w", err)
				}
				esiData.CorporationInfos[attacker.CorporationID] = *corpData
			}
		}

		if attacker.CharacterID != 0 {
			if _, exists := esiData.CharacterInfos[attacker.CharacterID]; !exists {
				charData, err := esi.GetCharacterInfo(ctx, attacker.CharacterID)
				if err != nil {
					var notFoundError *model.NotFoundError
					if errors.As(err, &notFoundError) {
						// Handle the 404 case silently or with a warning, without logging as an error
						//  esi.Logger.Infof("Character %d not found; skipping\n", attacker.CharacterID)
						return nil
					}
					return fmt.Errorf("failed to fetch character details for attacker: %w", err)
				}
				esiData.CharacterInfos[attacker.CharacterID] = *charData
			}
		}
	}

	return nil
}
