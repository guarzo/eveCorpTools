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

	"github.com/guarzo/zkillanalytics/internal/config"
	"github.com/guarzo/zkillanalytics/internal/data"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
)

const (
	ESIDataStaleDuration = 48 * time.Hour
	MinESIDataSize       = 20 // 50b
)

// OrchestrateService coordinates data fetching, aggregation, and persistence.
type OrchestrateService struct {
	KillMailService *KillMailService
	ESIService      *EsiService
	InvTypeService  *data.InvTypeService
	Failed          *model.FailedCharacters
	Cache           *persist.Cache
	Logger          *logrus.Logger
	Client          *http.Client

	// Mutex to ensure only one GetAllData runs at a time
	mu              sync.Mutex
	mutexAcquiredAt time.Time // Tracks when the mutex was acquired
}

// NewOrchestrateService initializes and returns a new OrchestrateService instance.
func NewOrchestrateService(
	esiService *EsiService,
	killMailService *KillMailService,
	invTypeService *data.InvTypeService,
	failed *model.FailedCharacters,
	cache *persist.Cache,
	logger *logrus.Logger,
	client *http.Client,
) *OrchestrateService {
	return &OrchestrateService{
		ESIService:      esiService,
		KillMailService: killMailService,
		InvTypeService:  invTypeService,
		Failed:          failed,
		Cache:           cache,
		Logger:          logger,
		Client:          client,
	}
}

// GetAllData orchestrates the data fetching process based on availability and necessity.
// It ensures that only one instance runs at a time.
func (os *OrchestrateService) GetAllData(ctx context.Context, corporations, alliances, characters []int, startString, endString string) (*model.ChartData, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	os.Logger.Infof("Attempting to acquire mutex...")
	if !os.AcquireMutex(5 * time.Second) {
		return nil, fmt.Errorf("another GetAllData operation is in progress")
	}

	// Use defer to ensure the mutex is always released, even if a panic occurs.
	defer func() {
		if r := recover(); r != nil {
			os.Logger.Errorf("Recovered from panic: %v", r)
		}
		os.Logger.Infof("Releasing mutex")
		os.ReleaseMutex()
	}()

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

	// Inside GetAllData
	dataAvailability, err := CheckDataAvailability(int(startDate.Month()), int(endDate.Month()), year)
	if err != nil {
		os.Logger.Errorf("Error checking data availability: %v", err)
		return nil, err
	}

	// Log available months for debugging
	availableMonths := []int{}
	for month, available := range dataAvailability {
		if available {
			availableMonths = append(availableMonths, month)
		}
	}
	os.Logger.Infof("Data available for months: %v", availableMonths)

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
	hardCodedIDs := &model.Ids{
		CorporationIDs: corporations,
		AllianceIDs:    alliances,
		CharacterIDs:   characters,
	}

	fetchIDs, err := persist.LoadIdsFromFile()
	if err != nil || fetchIDs == nil || (fetchIDs.CorporationIDs == nil && fetchIDs.AllianceIDs == nil && fetchIDs.CharacterIDs == nil) {
		os.Logger.Warnf("Using new ID File: %v", err)
		fetchIDs = &model.Ids{ // Initialize fetchIDs as a pointer to avoid nil reference issues
			CorporationIDs: make([]int, 0),
			AllianceIDs:    make([]int, 0),
			CharacterIDs:   make([]int, 0),
		}
		os.Logger.Infof("Loaded IDs from file")
	}

	idChanged, newIDs, err := persist.CheckIfIdsChanged(hardCodedIDs)
	if err != nil {
		newIDs = hardCodedIDs
		os.Logger.Errorf("Error checking if IDs changes")
	}

	// Create parameters for data fetching
	params := model.NewParams(os.Client, fetchIDs.CorporationIDs, fetchIDs.AllianceIDs, fetchIDs.CharacterIDs, year, esiData, idChanged, newIDs)

	// Update fetchIDs with new IDs if there were changes
	if idChanged {
		fetchIDs.CorporationIDs = append(fetchIDs.CorporationIDs, newIDs.CorporationIDs...)
		fetchIDs.AllianceIDs = append(fetchIDs.AllianceIDs, newIDs.AllianceIDs...)
		fetchIDs.CharacterIDs = append(fetchIDs.CharacterIDs, newIDs.CharacterIDs...)
	}

	// Fetch missing data if necessary
	newData, err := os.GetMissingData(ctx, &params, dataAvailability)
	if err != nil {
		os.Logger.Errorf("Error fetching missing data: %v", err)
		return nil, err
	}

	// Aggregate existing store data into ChartData
	for month := int(startDate.Month()); month <= int(endDate.Month()); month++ {
		if dataAvailability[month] {
			// Load existing data from file
			fileName := persist.GenerateZkillFileName(year, month)
			monthlyKillMailData, err := persist.ReadKillMailsFromFile(fileName)
			if err != nil {
				os.Logger.Errorf("Error loading data from file %s: %v", fileName, err)
				continue
			}
			// Populate tracked characters in ESIData
			err = os.ESIService.LoadTrackedCharacters(ctx, monthlyKillMailData.KillMails, esiData)
			if err != nil {
				os.Logger.Errorf("Error loading tracked characters into ESI data: %v", err)
				return nil, err
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

	if saveErr := persist.SaveFailedCharacters(os.Failed); saveErr != nil {
		os.Logger.Errorf("Error saving IDs data: %v", err)
	}

	cacheFile := persist.GenerateCacheDataFileName()
	err = os.ESIService.Cache.SaveToFile(cacheFile)
	if err != nil {
		os.Logger.Errorf("Error saving cache: %v", err)
	}

	fetchTotalTime := time.Since(fetchStart)
	os.Logger.Infof("Data fetching complete in %.2f seconds", fetchTotalTime.Seconds())
	return chartData, nil
}

func (os *OrchestrateService) isESIDataStale(fileInfo fs.FileInfo) bool {
	return time.Since(fileInfo.ModTime()) > ESIDataStaleDuration || fileInfo.Size() <= MinESIDataSize
}

func (os *OrchestrateService) GetMissingData(ctx context.Context, params *model.Params, dataAvailability map[int]bool) (*model.KillMailData, error) {
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

		// Fetch the data for this month
		monthlyKillMailData, err := os.KillMailService.GetKillMailDataForMonth(ctx, params, month)
		if err != nil {
			os.Logger.Errorf("Error fetching data for month %d: %v", month, err)
			return nil, err
		}

		// Reset aggregatedData.KillMails for each month to avoid carryover
		aggregatedData.KillMails = []model.DetailedKillMail{}

		// Aggregate and prepare to save
		aggregatedData.KillMails = os.KillMailService.AggregateKillMailDumps(aggregatedData.KillMails, monthlyKillMailData.KillMails)

		// Save the aggregated data to a unique store file
		fileName := persist.GenerateZkillFileName(params.Year, month)
		os.Logger.Infof("Saving data for month %d to file %s with %d killmails", month, fileName, len(aggregatedData.KillMails))
		if err = persist.SaveKillMailsToFile(fileName, monthlyKillMailData); err != nil {
			os.Logger.Errorf("Failed to save fetched data to file %s: %v", fileName, err)
			return nil, fmt.Errorf("failed to save fetched data: %w", err)
		}
	}

	return aggregatedData, nil
}

// AcquireMutex attempts to acquire the mutex within the given timeout.
func (os *OrchestrateService) AcquireMutex(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		os.mu.Lock()
		os.mutexAcquiredAt = time.Now() // Record when the mutex was acquired
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (os *OrchestrateService) ReleaseMutex() {
	// Calculate how long the mutex was held
	duration := time.Since(os.mutexAcquiredAt)

	// Log the duration
	os.Logger.Infof("Mutex was held for: %v", duration)

	// Release the mutex
	os.mu.Unlock()
}

// GetTrackedCorporations returns the list of tracked corporation IDs.
func (os *OrchestrateService) GetTrackedCorporations() []int {
	return config.CorporationIDs
}

// GetTrackedAlliances returns the list of tracked alliance IDs.
func (os *OrchestrateService) GetTrackedAlliances() []int {
	return config.AllianceIDs
}

// GetTrackedCharacters returns the list of tracked character IDs.
func (os *OrchestrateService) GetTrackedCharacters() []int {
	return config.CharacterIDs
}

// GetTrackedCharactersFromKillMails extracts tracked character IDs from killmails and ESI data.
func (os *OrchestrateService) GetTrackedCharactersFromKillMails(fullKillMail []model.DetailedKillMail, esiData *model.ESIData) []int {
	var trackedCharacters []int

	os.Logger.Debugf("tracked characters, killmail length: %d", len(fullKillMail))

	for _, km := range fullKillMail {
		for _, attacker := range km.Attackers {
			if persist.Contains(trackedCharacters, attacker.CharacterID) {
				os.Logger.Debugf("Character %d already tracked, skipping", attacker.CharacterID)
				continue
			}

			// Verify esiData contains attacker.CharacterID
			_, exists := esiData.CharacterInfos[attacker.CharacterID]
			if !exists {
				continue
			}

			// Verify esiData contains attacker.CorporationID
			corpInfo, exists := esiData.CorporationInfos[attacker.CorporationID]
			if !exists {
				continue
			}

			allianceID := corpInfo.AllianceID

			// Check DisplayCharacter
			if config.DisplayCharacter(attacker.CharacterID, attacker.CorporationID, allianceID) {
				trackedCharacters = append(trackedCharacters, attacker.CharacterID)
			}
		}
	}

	os.Logger.Debugf("Found %d tracked characters", len(trackedCharacters))
	return trackedCharacters
}

func (os *OrchestrateService) LookupType(id int) string {
	return os.InvTypeService.QueryInvType(id)
}

// CheckDataAvailability checks which months within a range have data files already present.
func CheckDataAvailability(startMonth, endMonth, year int) (map[int]bool, error) {
	dataAvailability := make(map[int]bool)
	currentMonth := int(time.Now().Month())
	previousMonth := currentMonth - 1
	if previousMonth == 0 {
		previousMonth = 12 // Handle year change from January to December of previous year
		year -= 1
	}

	for month := startMonth; month <= endMonth; month++ {
		fileName := persist.GenerateZkillFileName(year, month)
		if fileInfo, err := fs.Stat(fileName); err == nil {
			if fileInfo.Size() <= 1*1024 {
				fmt.Printf("File %s is too small (%d bytes). Marking as unavailable.\n", fileName, fileInfo.Size())
				dataAvailability[month] = false
				continue
			}
			fmt.Printf("Data for %04d-%02d already exists.\n", year, month)
			dataAvailability[month] = true

			// Handle current month: force update if file is stale
			if month == currentMonth {
				fileDate := fileInfo.ModTime().Truncate(24 * time.Hour)
				today := time.Now().Truncate(24 * time.Hour)
				if fileDate.Before(today) {
					fmt.Printf("Removing stale month-to-date file %s...\n", fileName)
					dataAvailability[month] = false
					err = fs.Remove(fileName)
					if err != nil {
						fmt.Printf("Error removing stale month-to-date file %s: %s\n", fileName, err)
					}
				} else {
					fmt.Printf("Continuing to use current month data for %s, %s\n", fileName, fileDate.Format("2006-01-02:15:04:05"))
				}
			}

			// Handle previous month: always set as unavailable to force re-fetch
			if month == previousMonth {
				fileDate := fileInfo.ModTime().Truncate(24 * time.Hour)
				today := time.Now().Truncate(24 * time.Hour)
				if fileDate.Before(today) {
					fmt.Printf("Removing stale last month file %s...\n", fileName)
					dataAvailability[month] = false
					err = fs.Remove(fileName)
					if err != nil {
						fmt.Printf("Error removing stale last month %s: %s\n", fileName, err)
					}
				} else {
					fmt.Printf("Continuing to use existing last month data for %s, %s\n", fileName, fileDate.Format("2006-01-02:15:04:05"))
				}
			}
		} else {
			dataAvailability[month] = false
		}
	}
	return dataAvailability, nil
}
