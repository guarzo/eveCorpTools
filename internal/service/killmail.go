// internal/service/killmail_service.go

package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/gambtho/zkillanalytics/internal/api/zkill"
	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
)

// KillMailService handles killmail-related operations.
type KillMailService struct {
	ZKillClient *zkill.ZkillClient
	EsiService  *EsiService
	Cache       *persist.Cache
	Logger      *logrus.Logger
}

// NewKillMailService creates a new instance of KillMailService.
func NewKillMailService(zkillClient *zkill.ZkillClient, esiService *EsiService, cache *persist.Cache, logger *logrus.Logger) *KillMailService {
	return &KillMailService{
		ZKillClient: zkillClient,
		EsiService:  esiService,
		Cache:       cache,
		Logger:      logger,
	}
}

// GetKillMailDataForMonth fetches and processes data for a specific month.
func (km *KillMailService) GetKillMailDataForMonth(ctx context.Context, params *config.Params, month int) (*model.KillMailData, error) {
	aggregatedMonthData := &model.KillMailData{
		KillMails: []model.DetailedKillMail{},
	}

	// Create a map to keep track of processed killmail IDs to avoid duplicates.
	killMailIDs := make(map[int]bool)

	// Define entity groups.
	entityGroups := map[string][]int{
		config.EntityTypeCorporation: params.Corporations,
		config.EntityTypeAlliance:    params.Alliances,
		config.EntityTypeCharacter:   params.Characters,
	}

	km.Logger.Infof("Fetching data for %04d-%02d...", params.Year, month)

	// Iterate over each entity type and their IDs to fetch killmails.
	for entityType, entityIDs := range entityGroups {
		for _, entityID := range entityIDs {
			// Fetch killmails for the entity.
			killMails, err := km.ZKillClient.GetKillsPageData(entityType, entityID, 1, params.Year, month)
			if err != nil {
				km.Logger.Errorf("Error fetching kills for %s ID %d: %v", entityType, entityID, err)
				continue
			}

			// Process fetched killmails.
			err = km.processKillMails(ctx, killMails, killMailIDs, aggregatedMonthData)
			if err != nil {
				km.Logger.Errorf("Error processing kills for %s ID %d: %v", entityType, entityID, err)
				continue
			}

			// Fetch loss killmails for the entity.
			lossKillMails, err := km.ZKillClient.GetLossPageData(entityType, entityID, 1, params.Year, month)
			if err != nil {
				km.Logger.Errorf("Error fetching losses for %s ID %d: %v", entityType, entityID, err)
				continue
			}

			// Process fetched loss killmails.
			err = km.processKillMails(ctx, lossKillMails, killMailIDs, aggregatedMonthData)
			if err != nil {
				km.Logger.Errorf("Error processing losses for %s ID %d: %v", entityType, entityID, err)
				continue
			}
		}
	}

	// Fetch victim killmails.
	err := km.GetVictimKillMails(ctx, params, month, aggregatedMonthData, killMailIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching victim kill mails: %w", err)
	}

	return aggregatedMonthData, nil
}

// processKillMails processes a slice of KillMail and updates aggregated data.
func (km *KillMailService) processKillMails(ctx context.Context, killMails []model.KillMail, killMailIDs map[int]bool, aggregatedData *model.KillMailData) error {
	for index, mail := range killMails {
		// Check if the killmail ID is already processed.
		if _, exists := killMailIDs[int(mail.KillMailID)]; exists {
			continue
		}

		// Log progress every 100 killmails.
		if index%100 == 0 {
			km.Logger.Infof("Processing killmail ID %d...%d of %d", mail.KillMailID, index, len(killMails))
		}

		// Process the full killmail and update aggregated data.
		err := km.AddEsiKillMail(ctx, mail, aggregatedData)
		if err != nil {
			km.Logger.Errorf("Error processing kill mail ID %d: %v", mail.KillMailID, err)
			continue
		}

		// Mark the killmail ID as processed.
		killMailIDs[int(mail.KillMailID)] = true
	}

	return nil
}

// GetVictimKillMails fetches killmails where the victim is part of specified entities.
func (km *KillMailService) GetVictimKillMails(ctx context.Context, params *config.Params, month int, aggregatedData *model.KillMailData, killMailIDs map[int]bool) error {
	// Define entity groups for victims.
	entityGroups := map[string][]int{
		config.EntityTypeCorporation: params.Corporations,
		config.EntityTypeAlliance:    params.Alliances,
		config.EntityTypeCharacter:   params.Characters,
	}

	km.Logger.Infof("Fetching victim kill mails for %04d-%02d...", params.Year, month)

	// Iterate over each entity type and their IDs.
	for entityType, entityIDs := range entityGroups {
		for _, entityID := range entityIDs {
			// Fetch victim killmails using the EsiClient.
			victimKillMails, err := km.ZKillClient.GetVictimKillsPageData(entityType, entityID, 1, params.Year, month)
			if err != nil {
				km.Logger.Errorf("Error fetching victim kills for %s ID %d: %v", entityType, entityID, err)
				continue
			}

			// Process fetched victim killmails.
			err = km.processKillMails(ctx, victimKillMails, killMailIDs, aggregatedData)
			if err != nil {
				km.Logger.Errorf("Error processing victim kills for %s ID %d: %v", entityType, entityID, err)
				continue
			}
		}
	}

	return nil
}

// AggregateKillMailDumps combines KillMailData into ChartData.
func (km *KillMailService) AggregateKillMailDumps(base, addition *model.KillMailData) *model.KillMailData {
	if base == nil {
		return addition
	}
	if addition == nil {
		return base
	}

	base.KillMails = append(base.KillMails, addition.KillMails...)
	return base
}

func (km *KillMailService) AddEsiKillMail(ctx context.Context, mail model.KillMail, aggregatedData *model.KillMailData) error {
	fullKillMail, err := km.EsiService.EsiClient.GetEsiKillMail(ctx, int(mail.KillMailID), mail.ZKB.Hash)
	if err != nil {
		return fmt.Errorf("failed to fetch full killmail for ID %d: %s", mail.KillMailID, err)
	}
	dKM := model.DetailedKillMail{KillMail: mail, EsiKillMail: *fullKillMail}
	aggregatedData.KillMails = append(aggregatedData.KillMails, dKM)
	return nil
}
