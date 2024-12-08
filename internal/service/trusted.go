package service

import (
	"fmt"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/sirupsen/logrus"
)

// TrustedService provides methods to manage trusted and untrusted entities.
type TrustedService struct {
	dataLoader func() (*model.TrustedCharacters, error)
	dataSaver  func(*model.TrustedCharacters) error
	Logger     *logrus.Logger
}

// NewTrustedService creates a new TrustedService with injected dependencies.
func NewTrustedService(dataLoader func() (*model.TrustedCharacters, error), dataSaver func(*model.TrustedCharacters) error, logger *logrus.Logger) *TrustedService {
	return &TrustedService{
		dataLoader: dataLoader,
		dataSaver:  dataSaver,
		Logger:     logger,
	}
}

// AddTrustedCharacter adds a new character to the trusted list.
func (s *TrustedService) AddTrustedCharacter(newCharacter model.TrustedCharacter) error {
	s.Logger.Infof("Starting addition process for character ID: %d (%s)", newCharacter.CharacterID, newCharacter.CharacterName)

	trustedData, err := s.dataLoader()
	if err != nil {
		s.Logger.Errorf("Failed to load trusted data during addition: %v", err)
		return fmt.Errorf("failed to load trusted data: %v", err)
	}

	for _, char := range trustedData.TrustedCharacters {
		if char.CharacterID == newCharacter.CharacterID {
			s.Logger.Infof("Character ID %d (%s) is already in the trusted list", newCharacter.CharacterID, newCharacter.CharacterName)
			return nil
		}
	}

	trustedData.TrustedCharacters = append(trustedData.TrustedCharacters, newCharacter)
	s.Logger.Infof("Character ID %d (%s) added to the trusted list", newCharacter.CharacterID, newCharacter.CharacterName)

	err = s.dataSaver(trustedData)
	if err != nil {
		s.Logger.Errorf("Failed to save updated trusted data after addition: %v", err)
		return fmt.Errorf("failed to save updated trusted data: %v", err)
	}

	s.Logger.Infof("Trusted data successfully saved after addition. Total trusted characters: %d", len(trustedData.TrustedCharacters))
	return nil
}

// RemoveTrustedCharacter removes a character from the trusted list by CharacterID.
func (s *TrustedService) RemoveTrustedCharacter(characterID int64) error {
	s.Logger.Infof("Starting removal process for character ID: %d", characterID)

	trustedData, err := s.dataLoader()
	if err != nil {
		s.Logger.Errorf("Failed to load trusted data during removal: %v", err)
		return fmt.Errorf("failed to load trusted data: %v", err)
	}

	initialCount := len(trustedData.TrustedCharacters)
	trustedData.TrustedCharacters = filterCharacters(trustedData.TrustedCharacters, characterID)
	finalCount := len(trustedData.TrustedCharacters)

	if initialCount == finalCount {
		s.Logger.Warnf("Character ID %d was not found in the trusted list", characterID)
	} else {
		s.Logger.Infof("Character ID %d successfully removed from the trusted list", characterID)
	}

	err = s.dataSaver(trustedData)
	if err != nil {
		s.Logger.Errorf("Failed to save updated trusted data after removal: %v", err)
		return fmt.Errorf("failed to save updated trusted data: %v", err)
	}

	s.Logger.Infof("Trusted data successfully saved after removal. Total remaining characters: %d", finalCount)
	return nil
}

// AddTrustedCorporation adds a new corporation to the trusted list.
func (s *TrustedService) AddTrustedCorporation(newCorporation model.TrustedCorporation) error {
	trustedData, err := s.dataLoader()
	if err != nil {
		return fmt.Errorf("failed to load trusted data: %v", err)
	}

	for _, corp := range trustedData.TrustedCorporations {
		if corp.CorporationID == newCorporation.CorporationID {
			s.Logger.Infof("Corporation %d already trusted", newCorporation.CorporationID)
			return nil
		}
	}

	trustedData.TrustedCorporations = append(trustedData.TrustedCorporations, newCorporation)
	return s.dataSaver(trustedData)
}

// RemoveTrustedCorporation removes a corporation from the trusted list by CorporationID.
func (s *TrustedService) RemoveTrustedCorporation(id int64) error {
	trustedData, err := s.dataLoader()
	if err != nil {
		return fmt.Errorf("failed to load trusted data: %v", err)
	}

	trustedData.TrustedCorporations = filterCorporations(trustedData.TrustedCorporations, id)
	s.Logger.Infof("Removed corporation %d from trusted list", id)
	return s.dataSaver(trustedData)
}

// AddUntrustedCharacter adds a character to the untrusted list.
func (s *TrustedService) AddUntrustedCharacter(character model.TrustedCharacter) error {
	data, err := s.dataLoader()
	if err != nil {
		return fmt.Errorf("failed to load trusted data: %v", err)
	}

	for _, existing := range data.UntrustedCharacters {
		if existing.CharacterID == character.CharacterID {
			s.Logger.Infof("Character %d already in untrusted list", character.CharacterID)
			return nil
		}
	}

	data.UntrustedCharacters = append(data.UntrustedCharacters, character)
	return s.dataSaver(data)
}

// AddUntrustedCorporation adds a corporation to the untrusted list.
func (s *TrustedService) AddUntrustedCorporation(corp model.TrustedCorporation) error {
	data, err := s.dataLoader()
	if err != nil {
		return fmt.Errorf("failed to load trusted data: %v", err)
	}

	for _, existing := range data.UntrustedCorporations {
		if existing.CorporationID == corp.CorporationID {
			s.Logger.Infof("Corporation %d already in untrusted list", corp.CorporationID)
			return nil
		}
	}

	data.UntrustedCorporations = append(data.UntrustedCorporations, corp)
	return s.dataSaver(data)
}

// RemoveUntrustedCharacter removes a character from the untrusted list by CharacterID.
func (s *TrustedService) RemoveUntrustedCharacter(characterID int64) error {
	data, err := s.dataLoader()
	if err != nil {
		return fmt.Errorf("failed to load trusted data: %v", err)
	}

	var filtered []model.TrustedCharacter
	for _, char := range data.UntrustedCharacters {
		if char.CharacterID != characterID {
			filtered = append(filtered, char)
		}
	}

	if len(filtered) == len(data.UntrustedCharacters) {
		s.Logger.Infof("Character %d not found in untrusted list", characterID)
		return nil
	}

	data.UntrustedCharacters = filtered
	return s.dataSaver(data)
}

// RemoveUntrustedCorporation removes a corporation from the untrusted list by CorporationID.
func (s *TrustedService) RemoveUntrustedCorporation(corpID int64) error {
	data, err := s.dataLoader()
	if err != nil {
		return fmt.Errorf("failed to load trusted data: %v", err)
	}

	var filtered []model.TrustedCorporation
	for _, corp := range data.UntrustedCorporations {
		if corp.CorporationID != corpID {
			filtered = append(filtered, corp)
		}
	}

	if len(filtered) == len(data.UntrustedCorporations) {
		s.Logger.Infof("Corporation %d not found in untrusted list", corpID)
		return nil
	}

	data.UntrustedCorporations = filtered
	return s.dataSaver(data)
}

func (s *TrustedService) GetTrustedCharacters() (*model.TrustedCharacters, error) {
	return s.dataLoader()
}

// Utility function to filter out a character by ID.
func filterCharacters(characters []model.TrustedCharacter, excludeID int64) []model.TrustedCharacter {
	updated := make([]model.TrustedCharacter, 0, len(characters))
	for _, char := range characters {
		if char.CharacterID != excludeID {
			updated = append(updated, char)
		}
	}
	return updated
}

// Utility function to filter out a corporation by ID.
func filterCorporations(corporations []model.TrustedCorporation, excludeID int64) []model.TrustedCorporation {
	updated := make([]model.TrustedCorporation, 0, len(corporations))
	for _, corp := range corporations {
		if corp.CorporationID != excludeID {
			updated = append(updated, corp)
		}
	}
	return updated
}
