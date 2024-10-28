// internal/service/killmail_service.go

package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/zkillanalytics/internal/api/zkill"
	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
)

// KillMailService handles killmail-related operations.
type KillMailService struct {
	ZKillClient *zkill.ZkillClient
	EsiService  *EsiService
	Cache       *persist.Cache
	Logger      *logrus.Logger

	// Mutex to protect killMailIDs map if accessed concurrently
	killMailMu sync.Mutex
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

	const maxPages = 100 // Define a sensible maximum number of pages
	processedKillMails := 0

	for entityType, entityIDs := range entityGroups {
		for _, entityID := range entityIDs {
			// Fetch kills
			page := 1
			for page <= maxPages {
				killMails, err := km.ZKillClient.GetKillsPageData(ctx, entityType, entityID, page, params.Year, month)
				if err != nil {
					km.Logger.Errorf("Error fetching kills for %s ID %d page %d: %v", entityType, entityID, page, err)
					break
				}
				if len(killMails) == 0 {
					break
				}

				err = km.processKillMails(ctx, killMails, killMailIDs, aggregatedMonthData)
				if err != nil {
					break
				}

				page++
				processedKillMails += len(killMails)
			}

			// Fetch losses
			page = 1
			for page <= maxPages {
				lossKillMails, err := km.ZKillClient.GetLossPageData(ctx, entityType, entityID, page, params.Year, month)
				if err != nil {
					km.Logger.Errorf("Error fetching losses for %s ID %d page %d: %v", entityType, entityID, page, err)
					break
				}

				if len(lossKillMails) == 0 {
					break
				}

				err = km.processKillMails(ctx, lossKillMails, killMailIDs, aggregatedMonthData)
				if err != nil {
					km.Logger.Errorf("Error processing losses for %s ID %d page %d: %v", entityType, entityID, page, err)
					break
				}

				page++
				processedKillMails += len(lossKillMails)
			}
		}
	}

	return aggregatedMonthData, nil
}

// processKillMails processes a slice of KillMail and updates aggregated data.
func (km *KillMailService) processKillMails(ctx context.Context, killMails []model.KillMail, killMailIDs map[int]bool, aggregatedData *model.KillMailData) error {
	for index, mail := range killMails {
		// Check if the killmail ID is already processed.
		if _, exists := killMailIDs[int(mail.KillMailID)]; exists {
			km.Logger.Debugf("Killmail ID %d already processed. Skipping.", mail.KillMailID)
			continue
		}

		// Process the full killmail and update aggregated data.
		err := km.AddEsiKillMail(ctx, mail, aggregatedData)
		if err != nil {
			km.Logger.Errorf("Error processing kill mail ID %d: %v", mail.KillMailID, err)
			continue
		}

		// Mark the killmail ID as processed.
		killMailIDs[int(mail.KillMailID)] = true

		// Optional: Log progress at intervals
		if (index+1)%100 == 0 {
			km.Logger.Infof("Processed %d killmails out of %d", index+1, len(killMails))
		}
	}
	km.Logger.Debugf("Processed %d killmails", len(killMails))

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
