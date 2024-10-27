// internal/service/orchestrate_service.go

package service

import (
	"context"
	"fmt"
	"net/http"
	fs "os"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/data"
	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/utils"
)

// OrchestrateService coordinates data fetching, aggregation, and persistence.
type OrchestrateService struct {
	KillMailService *KillMailService
	ESIService      *EsiService
	InvTypeService  *data.InvTypeService
	Cache           *persist.Cache
	Logger          *logrus.Logger
	Client          *http.Client

	// Atomic flag to ensure only one GetAllData runs at a time
	isRunning uint32
}

// NewOrchestrateService initializes and returns a new OrchestrateService instance.
func NewOrchestrateService(
	esiService *EsiService,
	killMailService *KillMailService,
	invTypeService *data.InvTypeService,
	cache *persist.Cache,
	logger *logrus.Logger,
	client *http.Client,
) *OrchestrateService {
	return &OrchestrateService{
		ESIService:      esiService,
		KillMailService: killMailService,
		InvTypeService:  invTypeService,
		Cache:           cache,
		Logger:          logger,
		Client:          client,
	}
}

// GetAllData orchestrates the data fetching process based on availability and necessity.
// It ensures that only one instance runs at a time.
func (os *OrchestrateService) GetAllData(ctx context.Context, corporations, alliances, characters []int, startString, endString string) (*model.ChartData, error) {
	// Attempt to acquire the lock using atomic flag
	if !atomic.CompareAndSwapUint32(&os.isRunning, 0, 1) {
		return nil, fmt.Errorf("another GetAllData operation is in progress")
	}
	// Ensure the flag is reset after the operation
	defer atomic.StoreUint32(&os.isRunning, 0)

	// Parse the start and end dates
	startDate, err := time.Parse("2006-01-02", startString)
	if err != nil {
		os.Logger.Errorf("Invalid start date format: %v", err)
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", endString)
	if err != nil {
		os.Logger.Errorf("Invalid end date format: %v", err)
		return nil, err
	}

	os.Logger.Infof("Fetching data from %s to %s...", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	fetchStart := time.Now()

	year := startDate.Year()
	esiRefresh := false

	// Check data availability
	dataAvailability, err := utils.CheckDataAvailability(int(startDate.Month()), int(endDate.Month()), year)
	if err != nil {
		os.Logger.Errorf("Error checking data availability: %v", err)
		return nil, err
	}

	// Load ESI data

	esiFileName := persist.GenerateEsiDataFileName()
	fileInfo, err := fs.Stat(esiFileName)

	esiData, err := persist.ReadEsiDataFromFile(esiFileName)
	if err != nil || os.isESIDataStale(fileInfo) {
		fmt.Println("Error loading esi data from file:", err)
		esiData = &model.ESIData{
			AllianceInfos:    make(map[int]model.Alliance),
			CharacterInfos:   make(map[int]model.Character),
			CorporationInfos: make(map[int]model.Corporation),
		}
		esiRefresh = true
	}

	// Check if IDs have changed
	fetchIDs := &model.Ids{
		CorporationIDs: corporations,
		AllianceIDs:    alliances,
		CharacterIDs:   characters,
	}

	idChanged, newIDs, err := persist.CheckIfIdsChanged(fetchIDs)
	if err != nil {
		os.Logger.Errorf("Error checking if IDs changed: %v", err)
		return nil, err
	}

	// Create parameters for data fetching
	params := config.NewParams(os.Client, corporations, alliances, characters, year, esiData, idChanged, newIDs)

	// Fetch missing data if necessary
	newData, err := os.GetMissingData(ctx, &params, dataAvailability)
	if err != nil {
		os.Logger.Errorf("Error fetching missing data: %v", err)
		return nil, err
	}

	// Aggregate existing monthly data into ChartData
	for month := int(startDate.Month()); month <= int(endDate.Month()); month++ {
		if dataAvailability[month] {
			// Load existing data from file
			fileName := persist.GenerateZkillFileName(year, month)
			monthlyKillMailData, err := persist.ReadKillMailDataFromFile(fileName)
			if err != nil {
				os.Logger.Errorf("Error loading data from file %s: %v", fileName, err)
				continue
			}
			// Aggregate KillMailData into ChartData
			newData = os.KillMailService.AggregateKillMailDumps(newData, monthlyKillMailData)
		}
	}

	// Initialize ChartData if it's nil
	chartData := &model.ChartData{
		KillMails: newData.KillMails,
		ESIData:   *esiData,
	}

	// Refresh ESI data if necessary
	if esiRefresh {
		err = os.ESIService.RefreshCharacter(ctx, chartData, os.Client)
		if err != nil {
			os.Logger.Errorf("Error refreshing ESI data: %v", err)
			return nil, err
		}
	}

	// Persist ESI data and IDs
	err = persist.SaveEsiDataToFile(esiFileName, esiData)
	if err != nil {
		os.Logger.Errorf("Error saving ESI data: %v", err)
		return nil, err
	}

	err = persist.SaveIdsToFile(fetchIDs)
	if err != nil {
		os.Logger.Errorf("Error saving IDs data: %v", err)
		return nil, err
	}

	fetchTotalTime := time.Since(fetchStart)
	os.Logger.Infof("Data fetching complete in %.2f seconds", fetchTotalTime.Seconds())
	return chartData, nil
}

func (os *OrchestrateService) isESIDataStale(fileInfo fs.FileInfo) bool {
	return time.Since(fileInfo.ModTime()) > 48*time.Hour || fileInfo.Size() <= 50*1024
}

// GetMissingData fetches missing killmails based on data availability and parameters.
func (os *OrchestrateService) GetMissingData(ctx context.Context, params *config.Params, dataAvailability map[int]bool) (*model.KillMailData, error) {
	aggregatedData := &model.KillMailData{
		KillMails: []model.DetailedKillMail{},
	}

	for month, available := range dataAvailability {
		if available && !params.ChangedIDs {
			// Data for this month is already available and IDs haven't changed; skip fetching.
			continue
		}

		if params.ChangedIDs {
			// IDs have changed; perform a full data pull for this month.
			available = false
		}

		monthlyKillMailData, err := os.KillMailService.GetKillMailDataForMonth(ctx, params, month)
		if err != nil {
			os.Logger.Errorf("Error fetching data for month %d: %v", month, err)
			return nil, err
		}

		fileName := persist.GenerateZkillFileName(params.Year, month)
		if err = persist.SaveKillMailsToFile(fileName, monthlyKillMailData); err != nil {
			os.Logger.Errorf("Failed to save fetched data to file %s: %v", fileName, err)
			return nil, fmt.Errorf("failed to save fetched data: %w", err)
		}

		aggregatedData = os.KillMailService.AggregateKillMailDumps(aggregatedData, monthlyKillMailData)
	}

	return aggregatedData, nil
}

// CombineEsiAndKillMail processes a kill mail and updates the aggregated data.
func (os *OrchestrateService) CombineEsiAndKillMail(ctx context.Context, kmModel model.KillMail, aggregatedData *model.KillMailData, esiData *model.ESIData) error {
	// Fetch the full killmail details using EsiService.
	fullKillMail, err := os.ESIService.EsiClient.GetEsiKillMail(ctx, int(kmModel.KillMailID), kmModel.ZKB.Hash)
	if err != nil {
		return fmt.Errorf("failed to fetch full killmail for ID %d: %w", kmModel.KillMailID, err)
	}

	// Aggregate ESI data into the global ESIData.
	err = os.ESIService.AggregateEsiData(ctx, fullKillMail, esiData)
	if err != nil {
		return fmt.Errorf("failed to aggregate ESI data: %w", err)
	}

	// Create a DetailedKillMail and append it to the aggregated killmails.
	dKM := model.DetailedKillMail{
		KillMail:    kmModel,
		EsiKillMail: *fullKillMail,
	}
	aggregatedData.KillMails = append(aggregatedData.KillMails, dKM)

	return nil
}
