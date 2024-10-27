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

func (km *KillMailService) GetKillMailDataForMonth(ctx context.Context, params *model.Params, month int) (*model.KillMailData, error) {
	aggregatedMonthData := &model.KillMailData{
		KillMails: []model.DetailedKillMail{},
	}

	killMailIDs := make(map[int]bool)

	entityGroups := map[string][]int{
		config.EntityTypeCorporation: params.Corporations,
		config.EntityTypeAlliance:    params.Alliances,
		config.EntityTypeCharacter:   params.Characters,
	}

	km.Logger.Infof("Fetching data for %04d-%02d...", params.Year, month)

	for entityType, entityIDs := range entityGroups {
		for _, entityID := range entityIDs {
			page := 1
			for {
				killMails, err := km.ZKillClient.GetKillsPageData(entityType, entityID, page, params.Year, month)
				if err != nil {
					km.Logger.Errorf("Error fetching kills for %s ID %d page %d: %v", entityType, entityID, page, err)
					break
				}

				km.Logger.Debugf("for page %d, %d kills", page, len(killMails))
				if len(killMails) == 0 {
					break
				}

				err = km.processKillMails(ctx, killMails, killMailIDs, aggregatedMonthData)
				if err != nil {
					km.Logger.Errorf("Error processing kills for %s ID %d page %d: %v", entityType, entityID, page, err)
					break
				}
				km.Logger.Debugf("in kills on page %d", page)
				page++
			}

			// Repeat similarly for loss killmails
			page = 1
			for {
				lossKillMails, err := km.ZKillClient.GetLossPageData(entityType, entityID, page, params.Year, month)
				if err != nil {
					km.Logger.Errorf("Error fetching losses for %s ID %d page %d: %v", entityType, entityID, page, err)
					break
				}
				km.Logger.Debugf("for page %d, %d losses", page, len(lossKillMails))

				if len(lossKillMails) == 0 {
					break
				}

				err = km.processKillMails(ctx, lossKillMails, killMailIDs, aggregatedMonthData)
				if err != nil {
					km.Logger.Errorf("Error processing losses for %s ID %d page %d: %v", entityType, entityID, page, err)
					break
				}
				km.Logger.Debugf("in loss on page %d", page)

				page++
				km.Logger.Debugf("finished page %d", page)
			}
		}
	}

	km.Logger.Debugf("about to get victim kill mails")
	err := km.GetVictimKillMails(ctx, params, month, aggregatedMonthData, killMailIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching victim kill mails: %w", err)
	}

	return aggregatedMonthData, nil
}

// processKillMails processes a slice of KillMail and updates aggregated data.
func (km *KillMailService) processKillMails(ctx context.Context, killMails []model.KillMail, killMailIDs map[int]bool, aggregatedData *model.KillMailData) error {
	km.Logger.Debugf("%d killmails to process", len(killMails))
	for index, mail := range killMails {
		km.Logger.Debugf("Processing %d", index)
		// Check if the killmail ID is already processed.
		if _, exists := killMailIDs[int(mail.KillMailID)]; exists {
			continue
		}

		// Log progress every 100 killmails.
		//if index%100 == 0 {
		km.Logger.Infof("Processing killmail ID %d...%d of %d", mail.KillMailID, index, len(killMails))
		//}

		// Process the full killmail and update aggregated data.
		err := km.AddEsiKillMail(ctx, mail, aggregatedData)
		if err != nil {
			km.Logger.Errorf("Error processing kill mail ID %d: %v", mail.KillMailID, err)
			continue
		}

		// Mark the killmail ID as processed.
		killMailIDs[int(mail.KillMailID)] = true
	}

	km.Logger.Debugf("finished processKillMails")
	return nil
}

// GetVictimKillMails fetches killmails where the victim is part of specified entities.
func (km *KillMailService) GetVictimKillMails(ctx context.Context, params *model.Params, month int, aggregatedData *model.KillMailData, killMailIDs map[int]bool) error {
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
			km.Logger.Infof("victim kills probably needs paging....")
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

func (km *KillMailService) AggregateKillMailDumps(base, addition []model.DetailedKillMail) []model.DetailedKillMail {
	if base == nil {
		return addition
	}
	if addition == nil {
		return base
	}

	km.Logger.Debug("finished Aggregate KM Dumps")
	return append(base, addition...)
}

func (km *KillMailService) AddEsiKillMail(ctx context.Context, mail model.KillMail, aggregatedData *model.KillMailData) error {
	fullKillMail, err := km.EsiService.EsiClient.GetEsiKillMail(ctx, int(mail.KillMailID), mail.ZKB.Hash)
	if err != nil {
		return fmt.Errorf("failed to fetch full killmail for ID %d: %s", mail.KillMailID, err)
	}
	dKM := model.DetailedKillMail{KillMail: mail, EsiKillMail: *fullKillMail}
	aggregatedData.KillMails = append(aggregatedData.KillMails, dKM)
	km.Logger.Debugf("Added ESI to killmail %d", mail.KillMailID)
	return nil
}
