// internal/service/orchestrate_service.go

package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
func (svc *OrchestrateService) GetAllData(ctx context.Context, corporations, alliances, characters []int, startString, endString string) (*model.ChartData, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	svc.Logger.Infof("Attempting to acquire mutex...")
	if !svc.AcquireMutex() {
		return nil, fmt.Errorf("another GetAllData operation is in progress")
	}

	// Use defer to ensure the mutex is always released, even if a panic occurs.
	defer func() {
		if r := recover(); r != nil {
			svc.Logger.Errorf("Recovered from panic: %v", r)
		}
		svc.Logger.Infof("Releasing mutex")
		svc.ReleaseMutex()
	}()

	// Parse the start and end dates
	startDate, err := time.Parse("2006-01-02", startString)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", endString)
	if err != nil {
		return nil, fmt.Errorf("invalid e date format: %w", err)
	}

	svc.Logger.Infof("Fetching data from %s to %s...", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	fetchStart := time.Now()
	year := startDate.Year()
	esiRefresh := false

	// Inside GetAllData
	dataAvailability, err := svc.CheckDataAvailability(int(startDate.Month()), int(endDate.Month()), year)
	if err != nil {
		svc.Logger.Errorf("Error checking data availability: %v", err)
		return nil, err
	}

	// Log available months for debugging
	var availableMonths []int
	for key, available := range dataAvailability {
		if available {
			availableMonths = append(availableMonths, key)
		}
	}

	for _, key := range availableMonths {
		aYear, month := extractYearMonthKey(key)
		svc.Logger.Infof("Data available for %04d-%02d", aYear, month)
	}

	// Load ESI data
	esiFileName := persist.GenerateEsiDataFileName()
	fileInfo, err := os.Stat(esiFileName)

	esiData, err := persist.ReadEsiDataFromFile(esiFileName)
	if err != nil || svc.isESIDataStale(fileInfo) {
		svc.Logger.Warnf("Using new ESI File: %v", err)
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
		svc.Logger.Warnf("Using new ID File: %v", err)
		fetchIDs = &model.Ids{ // Initialize fetchIDs as a pointer to avoid nil reference issues
			CorporationIDs: make([]int, 0),
			AllianceIDs:    make([]int, 0),
			CharacterIDs:   make([]int, 0),
		}
		svc.Logger.Infof("Loaded IDs from file")
	}

	idChanged, newIDs, err := persist.CheckIfIdsChanged(hardCodedIDs)
	if err != nil {
		newIDs = hardCodedIDs
		svc.Logger.Errorf("Error checking if IDs changes")
	}

	// Create parameters for data fetching
	params := model.NewParams(svc.Client, fetchIDs.CorporationIDs, fetchIDs.AllianceIDs, fetchIDs.CharacterIDs, year, esiData, idChanged, newIDs)

	// Update fetchIDs with new IDs if there were changes
	if idChanged {
		fetchIDs.CorporationIDs = append(fetchIDs.CorporationIDs, newIDs.CorporationIDs...)
		fetchIDs.AllianceIDs = append(fetchIDs.AllianceIDs, newIDs.AllianceIDs...)
		fetchIDs.CharacterIDs = append(fetchIDs.CharacterIDs, newIDs.CharacterIDs...)
	}

	// Fetch missing data if necessary
	newData, err := svc.GetMissingData(ctx, &params, dataAvailability)
	if err != nil {
		svc.Logger.Errorf("Error fetching missing data: %v", err)
		return nil, err
	}

	yearMonths, err := generateYearMonthPairs(int(startDate.Month()), int(endDate.Month()), startDate.Year())
	if err != nil {
		svc.Logger.Errorf("Error generating year-month pairs: %v", err)
		return nil, err
	}

	for _, ym := range yearMonths {
		y, m := ym.Year, ym.Month
		key := getYearMonthKey(y, m)
		if dataAvailability[key] {
			// Load existing data from file
			fileName := persist.GenerateZkillFileName(y, m)
			monthlyKillMailData, err := persist.ReadKillMailsFromFile(fileName)
			if err != nil {
				svc.Logger.Errorf("Error loading data from file %s: %v", fileName, err)
				continue
			}
			// Populate tracked characters in ESIData
			err = svc.ESIService.LoadTrackedCharacters(ctx, monthlyKillMailData.KillMails, esiData)
			if err != nil {
				svc.Logger.Errorf("Error loading tracked characters into ESI data: %v", err)
				return nil, err
			}

			// Aggregate KillMailData into NewData
			newData.KillMails = svc.KillMailService.AggregateKillMailDumps(newData.KillMails, monthlyKillMailData.KillMails)
		}
	}

	// Initialize ChartData
	chartData := &model.ChartData{
		KillMails: newData.KillMails,
		ESIData:   *esiData,
	}

	// Refresh ESI data if necessary
	if esiRefresh {
		err = svc.ESIService.RefreshEsiData(ctx, chartData, svc.Client)
		if err != nil {
			svc.Logger.Errorf("Error refreshing ESI data: %v", err)
			return nil, err
		}
	}

	// Persist ESI data and IDs
	err = persist.SaveEsiDataToFile(esiFileName, esiData)
	if err != nil {
		svc.Logger.Errorf("Error saving ESI data: %v", err)
		return nil, err
	}

	err = persist.SaveIdsToFile(fetchIDs)
	if err != nil {
		svc.Logger.Errorf("Error saving IDs data: %v", err)
		return nil, err
	}

	if saveErr := persist.SaveFailedCharacters(svc.Failed); saveErr != nil {
		svc.Logger.Errorf("Error saving IDs data: %v", err)
	}

	cacheFile := persist.GenerateCacheDataFileName()
	err = svc.ESIService.Cache.SaveToFile(cacheFile)
	if err != nil {
		svc.Logger.Errorf("Error saving cache: %v", err)
	}

	fetchTotalTime := time.Since(fetchStart)
	svc.Logger.Infof("Data fetching complete in %.2f seconds", fetchTotalTime.Seconds())
	return chartData, nil
}

func (svc *OrchestrateService) isESIDataStale(fileInfo os.FileInfo) bool {
	return time.Since(fileInfo.ModTime()) > ESIDataStaleDuration || fileInfo.Size() <= MinESIDataSize
}

func (svc *OrchestrateService) GetMissingData(ctx context.Context, params *model.Params, dataAvailability map[int]bool) (*model.KillMailData, error) {
	aggregatedData := &model.KillMailData{
		KillMails: []model.DetailedKillMail{},
	}

	for key, available := range dataAvailability {
		if available && !params.ChangedIDs {
			// Data for this month is already available and IDs haven't changed; skip fetching.
			continue
		}

		if params.ChangedIDs {
			// IDs have changed; perform a full data pull for this month.
			available = false
		}

		// Extract year and month from key
		year, month := extractYearMonthKey(key)

		// Fetch the data for this month
		monthlyKillMailData, err := svc.KillMailService.GetKillMailDataForMonth(ctx, params, month)
		if err != nil {
			svc.Logger.Errorf("Error fetching data for %04d-%02d: %v", year, month, err)
			return nil, err
		}

		// Aggregate and prepare to save
		aggregatedData.KillMails = svc.KillMailService.AggregateKillMailDumps(aggregatedData.KillMails, monthlyKillMailData.KillMails)

		// Save the aggregated data to a unique store file
		fileName := persist.GenerateZkillFileName(year, month)
		svc.Logger.Infof("Saving data for %04d-%02d to file %s with %d killmails", year, month, fileName, len(aggregatedData.KillMails))
		if err = persist.SaveKillMailsToFile(fileName, monthlyKillMailData); err != nil {
			svc.Logger.Errorf("Failed to save fetched data to file %s: %v", fileName, err)
			return nil, fmt.Errorf("failed to save fetched data: %w", err)
		}
	}

	return aggregatedData, nil
}

func (svc *OrchestrateService) AcquireMutex() bool {
	if svc.mu.TryLock() {
		svc.mutexAcquiredAt = time.Now()
		return true
	} else {
		return false
	}
}

func (svc *OrchestrateService) ReleaseMutex() {
	// Calculate how long the mutex was held
	duration := time.Since(svc.mutexAcquiredAt)

	// Log the duration
	svc.Logger.Infof("Mutex was held for: %v", duration)

	// Release the mutex
	svc.mu.Unlock()
}

// GetTrackedCorporations returns the list of tracked corporation IDs.
func (svc *OrchestrateService) GetTrackedCorporations() []int {
	return config.CorporationIDs
}

// GetTrackedAlliances returns the list of tracked alliance IDs.
func (svc *OrchestrateService) GetTrackedAlliances() []int {
	return config.AllianceIDs
}

// GetTrackedCharacters returns the list of tracked character IDs.
func (svc *OrchestrateService) GetTrackedCharacters() []int {
	return config.CharacterIDs
}

// GetTrackedCharactersFromKillMails extracts tracked character IDs from killmails and ESI data.
func (svc *OrchestrateService) GetTrackedCharactersFromKillMails(fullKillMail []model.DetailedKillMail, esiData *model.ESIData) []int {
	var trackedCharacters []int

	svc.Logger.Debugf("tracked characters, killmail length: %d", len(fullKillMail))

	for _, km := range fullKillMail {
		for _, attacker := range km.Attackers {
			if persist.Contains(trackedCharacters, attacker.CharacterID) {
				svc.Logger.Debugf("Character %d already tracked, skipping", attacker.CharacterID)
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

	svc.Logger.Debugf("Found %d tracked characters", len(trackedCharacters))
	return trackedCharacters
}

func (svc *OrchestrateService) LookupType(id int) string {
	return svc.InvTypeService.QueryInvType(id)
}

type YearMonth struct {
	Year  int
	Month int
}

func (svc *OrchestrateService) CheckDataAvailability(startMonth, endMonth, startYear int) (map[int]bool, error) {
	dataAvailability := make(map[int]bool)
	currentTime := time.Now()
	stalenessDuration := 24 * time.Hour

	yearMonths, err := generateYearMonthPairs(startMonth, endMonth, startYear)
	if err != nil {
		return nil, err
	}

	for _, ym := range yearMonths {
		y, m := ym.Year, ym.Month
		key := getYearMonthKey(y, m)

		fileName := persist.GenerateZkillFileName(y, m)
		fileInfo, err := os.Stat(fileName)

		if err != nil {
			dataAvailability[key] = false
			continue
		}

		if fileInfo.Size() <= 1*1024 {
			svc.Logger.Warnf("File %s is too small (%d bytes). Marking as unavailable.\n", fileName, fileInfo.Size())
			dataAvailability[key] = false
			continue
		}

		dataAvailability[key] = true
		svc.Logger.Warnf("Data for %04d-%02d already exists.\n", y, m)

		// Check if the file is stale for current or previous month
		if isCurrentOrPreviousMonth(y, m, currentTime) {
			age := currentTime.Sub(fileInfo.ModTime())
			if age > stalenessDuration {
				svc.Logger.Warnf("Removing stale file %s (age: %v)...\n", fileName, age)
				dataAvailability[key] = false
				err = os.Remove(fileName)
				if err != nil {
					svc.Logger.Errorf("Error removing stale file %s: %s\n", fileName, err)
				}
			} else {
				svc.Logger.Infof("Using recent file %s (age: %v)\n", fileName, age)
			}
		}
	}

	return dataAvailability, nil
}

func generateYearMonthPairs(startMonth, endMonth, startYear int) ([]struct{ Year, Month int }, error) {
	var yearMonths []struct{ Year, Month int }

	startDate := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(startYear, time.Month(endMonth), 1, 0, 0, 0, 0, time.UTC)

	if endDate.Before(startDate) {
		// Handle crossing over to the next year
		endDate = endDate.AddDate(1, 0, 0)
	}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
		yearMonths = append(yearMonths, struct{ Year, Month int }{d.Year(), int(d.Month())})
	}

	return yearMonths, nil
}

// getYearMonthKey generates a unique integer key from a year and month.
// For example, year 2023 and month 7 become 202307.
func getYearMonthKey(year, month int) int {
	return year*100 + month
}

// extractYearMonthKey extracts the year and month from a key generated by getYearMonthKey.
func extractYearMonthKey(key int) (year int, month int) {
	year = key / 100
	month = key % 100
	return
}

func isCurrentOrPreviousMonth(year, month int, currentTime time.Time) bool {
	currentYear, currentMonth := currentTime.Year(), int(currentTime.Month())
	previousTime := currentTime.AddDate(0, -1, 0)
	previousYear, previousMonth := previousTime.Year(), int(previousTime.Month())

	return (year == currentYear && month == currentMonth) ||
		(year == previousYear && month == previousMonth)
}
