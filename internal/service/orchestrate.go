// internal/service/orchestrate_service.go

package service

import (
	"context"
	"fmt"
	"net/http"
	fs "os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gambtho/zkillanalytics/internal/config"
	"github.com/gambtho/zkillanalytics/internal/data"
	"github.com/gambtho/zkillanalytics/internal/model"
	"github.com/gambtho/zkillanalytics/internal/persist"
	"github.com/gambtho/zkillanalytics/internal/utils"
)

const (
	ESIDataStaleDuration = 48 * time.Hour
	MinESIDataSize       = 50 * 1024 // 50KB
)

// OrchestrateService coordinates data fetching, aggregation, and persistence.
type OrchestrateService struct {
	KillMailService *KillMailService
	ESIService      *EsiService
	InvTypeService  *data.InvTypeService
	Cache           *persist.Cache
	Logger          *logrus.Logger
	Client          *http.Client

	// Mutex to ensure only one GetAllData runs at a time
	mu sync.Mutex
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
	if !os.AcquireMutex(5 * time.Second) {
		return nil, fmt.Errorf("another GetAllData operation is in progress")
	}
	defer os.ReleaseMutex()

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
	time.Sleep(500000)
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
		os.Logger.Warnf("Using new ESI File: %v", err)
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

	idChanged, newIDs, idStr := persist.CheckIfIdsChanged(fetchIDs)
	if idStr != "" {
		os.Logger.Info("Checking if IDs changed: %v", idStr)
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
			// Aggregate KillMailData into NewData
			newData.KillMails = os.KillMailService.AggregateKillMailDumps(newData.KillMails, monthlyKillMailData.KillMails)
		}
	}

	// Initialize ChartData
	chartData := &model.ChartData{
		KillMails: newData.KillMails,
		ESIData:   *esiData,
	}

	// Refresh ESI data if necessary
	if esiRefresh {
		err = os.ESIService.RefreshEsiData(ctx, chartData, os.Client)
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
	return time.Since(fileInfo.ModTime()) > ESIDataStaleDuration || fileInfo.Size() <= MinESIDataSize
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

		aggregatedData.KillMails = os.KillMailService.AggregateKillMailDumps(aggregatedData.KillMails, monthlyKillMailData.KillMails)
	}

	return aggregatedData, nil
}

// AcquireMutex attempts to acquire the mutex within the given timeout.
func (os *OrchestrateService) AcquireMutex(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		os.mu.Lock()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// ReleaseMutex releases the mutex.
func (os *OrchestrateService) ReleaseMutex() {
	os.mu.Unlock()
}

// GetTrackedCorporations returns the list of tracked corporation IDs.
func (os *OrchestrateService) GetTrackedCorporations() []int {
	return persist.CorporationIDs
}

// GetTrackedAlliances returns the list of tracked alliance IDs.
func (os *OrchestrateService) GetTrackedAlliances() []int {
	return persist.AllianceIDs
}

// GetTrackedCharacters returns the list of tracked character IDs.
func (os *OrchestrateService) GetTrackedCharacters() []int {
	return persist.CharacterIDs
}

// GetTrackedCharactersFromKillMails extracts tracked character IDs from killmails and ESI data.
func (os *OrchestrateService) GetTrackedCharactersFromKillMails(fullKillMail []model.DetailedKillMail, esiData *model.ESIData) []int {
	var trackedCharacters []int

	for _, km := range fullKillMail {
		for _, attacker := range km.Attackers {
			if persist.Contains(trackedCharacters, attacker.CharacterID) {
				continue
			}

			_, exists := esiData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			corpInfo, exists := esiData.CorporationInfos[attacker.CorporationID]
			if !exists {
				continue
			}

			allianceID := corpInfo.AllianceID

			if persist.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, allianceID) {
				// fmt.Println(fmt.Sprintf("Adding character %d to tracked characters", attacker.CharacterID))
				trackedCharacters = append(trackedCharacters, attacker.CharacterID)
			}
		}
	}

	// fmt.Println(fmt.Sprintf("Found %d tracked characters", len(trackedCharacters)))
	return trackedCharacters
}

func (os *OrchestrateService) LookupType(id int) string {
	return os.InvTypeService.QueryInvType(id)
}
