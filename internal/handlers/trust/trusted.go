package trust

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/zkillanalytics/internal/handlers"
	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/persist"
	"github.com/guarzo/zkillanalytics/internal/xlog"

	"github.com/guarzo/zkillanalytics/internal/service"
)

// EntityData holds the resolved data for an entity, including corporation details if applicable.
type EntityData struct {
	ID              int64
	Name            string
	CorporationID   int64
	CorporationName string
	AllianceID      int64
	AllianceName    string
}

// Helper function to parse and resolve the identifier.
func resolveIdentifier(identifier string, entityType string, mainIdentity int64, token *oauth2.Token, esiService *service.EsiService) (EntityData, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return EntityData{}, fmt.Errorf("identifier is empty")
	}

	// Check if identifier is numeric.
	if id, err := strconv.ParseInt(identifier, 10, 64); err == nil {
		if id <= 0 {
			return EntityData{}, fmt.Errorf("identifier must be a positive number")
		}
		return EntityData{ID: id, Name: ""}, nil
	}

	// If identifier is not numeric, treat it as a name and resolve to ID.
	var resolvedID int32
	var err error
	if entityType == "character" {
		resolvedID, err = esiService.EsiClient.CharacterIDSearch(mainIdentity, identifier, token)
	} else if entityType == "corporation" {
		resolvedID, err = esiService.EsiClient.CorporationIDSearch(mainIdentity, identifier, token)
		if err != nil {
			esiService.Logger.Warnf("Failed to resolve %s identifier %s to an ID: %v", entityType, identifier, err)
			return EntityData{}, fmt.Errorf("identifier resolution failed: %v", err)
		}

	} else {
		return EntityData{}, fmt.Errorf("unknown entity type: %s", entityType)
	}

	if err != nil {
		return EntityData{}, fmt.Errorf("failed to resolve name to ID: %v", err)
	}

	if resolvedID <= 0 {
		return EntityData{}, fmt.Errorf("resolved ID is invalid for identifier: %s", identifier)
	}

	return EntityData{ID: int64(resolvedID), Name: identifier}, nil
}

func handleAddEntity(s *handlers.SessionService, trustedService *service.TrustedService, esiService *service.EsiService, w http.ResponseWriter, r *http.Request, trustStatus, entityType string) {
	var request struct {
		Identifier string `json:"identifier"`
	}

	setDirtyBit(s, r, w)

	// Decode the request body to access the identifier
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		trustedService.Logger.Warnf("Bad request payload: %v", err)
		handlers.WriteJSONError(w, "Invalid request payload", request.Identifier, http.StatusBadRequest, trustedService.Logger)
		return
	}

	// Retrieve session identity and token
	mainIdentity, token, err := handlers.GetSessionIdentity(s, r, trustedService.Logger)
	if err != nil {
		handlers.WriteJSONError(w, "Authentication required", request.Identifier, http.StatusUnauthorized, trustedService.Logger)
		return
	}

	// Resolve identifier
	resolvedData, err := resolveIdentifier(request.Identifier, entityType, mainIdentity, &token, esiService)
	if err != nil {
		trustedService.Logger.Warnf("Identifier resolution error: %v", err)
		handlers.WriteJSONError(w, "Identifier resolution failed", request.Identifier, http.StatusBadRequest, trustedService.Logger)
		return
	}

	// Fetch entity data after resolution
	fetchedData, err := fetchEntityData(entityType, resolvedData, &token, esiService)
	if err != nil {
		trustedService.Logger.Errorf("Entity data fetching error: %v", err)
		handlers.WriteJSONError(w, "Entity data retrieval failed", request.Identifier, http.StatusInternalServerError, trustedService.Logger)
		return
	}

	// Retrieve 'AddedBy' information
	addedBy, err := esiService.EsiClient.GetPublicCharacterData(mainIdentity, &token)
	if err != nil {
		trustedService.Logger.Errorf("Error retrieving character data for AddedBy: %v", err)
		handlers.WriteJSONError(w, "Failed to validate AddedBy character", request.Identifier, http.StatusInternalServerError, trustedService.Logger)
		return
	}

	// Prepare entity data for response
	var responseEntity interface{}
	switch {
	case trustStatus == "trusted" && entityType == "character":
		trustedCharacter := model.TrustedCharacter{
			CharacterID:     fetchedData.ID,
			CharacterName:   fetchedData.Name,
			CorporationID:   fetchedData.CorporationID,
			CorporationName: fetchedData.CorporationName,
			AddedBy:         addedBy.Name,
			DateAdded:       time.Now(),
		}
		if err := trustedService.AddTrustedCharacter(trustedCharacter); err != nil {
			trustedService.Logger.Errorf("Error saving trusted character: %v", err)
			handlers.WriteJSONError(w, "Failed to save trusted character", request.Identifier, http.StatusInternalServerError, trustedService.Logger)
			return
		}
		responseEntity = trustedCharacter

	case trustStatus == "trusted" && entityType == "corporation":
		trustedCorporation := model.TrustedCorporation{
			CorporationID:   fetchedData.ID,
			CorporationName: fetchedData.Name,
			AllianceID:      fetchedData.AllianceID,
			AllianceName:    fetchedData.AllianceName,
			AddedBy:         addedBy.Name,
			DateAdded:       time.Now(),
		}
		if err := trustedService.AddTrustedCorporation(trustedCorporation); err != nil {
			trustedService.Logger.Errorf("Error saving trusted corporation: %v", err)
			handlers.WriteJSONError(w, "Failed to save trusted corporation", request.Identifier, http.StatusInternalServerError, trustedService.Logger)
			return
		}
		responseEntity = trustedCorporation

	case trustStatus == "untrusted" && entityType == "character":
		untrustedCharacter := model.TrustedCharacter{
			CharacterID:     fetchedData.ID,
			CharacterName:   fetchedData.Name,
			CorporationID:   fetchedData.CorporationID,
			CorporationName: fetchedData.CorporationName,
			AddedBy:         addedBy.Name,
			DateAdded:       time.Now(),
		}
		if err := trustedService.AddUntrustedCharacter(untrustedCharacter); err != nil {
			trustedService.Logger.Errorf("Error saving untrusted character: %v", err)
			handlers.WriteJSONError(w, "Failed to save untrusted character", request.Identifier, http.StatusInternalServerError, trustedService.Logger)
			return
		}
		responseEntity = untrustedCharacter

	case trustStatus == "untrusted" && entityType == "corporation":
		untrustedCorporation := model.TrustedCorporation{
			CorporationID:   fetchedData.ID,
			CorporationName: fetchedData.Name,
			AllianceID:      fetchedData.AllianceID,
			AllianceName:    fetchedData.AllianceName,
			AddedBy:         addedBy.Name,
			DateAdded:       time.Now(),
		}
		if err := trustedService.AddUntrustedCorporation(untrustedCorporation); err != nil {
			trustedService.Logger.Errorf("Error saving untrusted corporation: %v", err)
			handlers.WriteJSONError(w, "Failed to save untrusted corporation", request.Identifier, http.StatusInternalServerError, trustedService.Logger)
			return
		}
		responseEntity = untrustedCorporation

	default:
		handlers.WriteJSONError(w, "Unsupported entity type or trust status", request.Identifier, http.StatusBadRequest, trustedService.Logger)
		return
	}

	// Send the full entity response back to the client
	handlers.WriteJSONResponse(w, responseEntity, http.StatusOK, trustedService.Logger)
}

func AddTrustedCharacterHandler(s *handlers.SessionService, trustedService *service.TrustedService, esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleAddEntity(s, trustedService, esiService, w, r, "trusted", "character")
	}
}

func AddTrustedCorporationHandler(s *handlers.SessionService, trustedService *service.TrustedService, esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleAddEntity(s, trustedService, esiService, w, r, "trusted", "corporation")
	}
}

func AddUntrustedCharacterHandler(s *handlers.SessionService, trustedService *service.TrustedService, esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleAddEntity(s, trustedService, esiService, w, r, "untrusted", "character")
	}
}

func AddUntrustedCorporationHandler(s *handlers.SessionService, trustedService *service.TrustedService, esiService *service.EsiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleAddEntity(s, trustedService, esiService, w, r, "untrusted", "corporation")
	}
}

func handleRemoveEntity(s *handlers.SessionService, trustedService *service.TrustedService, w http.ResponseWriter, r *http.Request, trustStatus, entityType string) {
	// Decode request body to retrieve 'identifier'
	var request struct {
		Identifier string `json:"identifier"`
	}

	setDirtyBit(s, r, w)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		trustedService.Logger.Warnf("Bad request payload: %v", err)
		handlers.WriteJSONResponse(w, handlers.ErrorResponse{Error: "Invalid request payload"}, http.StatusBadRequest, trustedService.Logger)
		return
	}

	trustedService.Logger.Infof("Removing %s %s with identifier: %v", trustStatus, entityType, request.Identifier)

	// Parse identifier as an integer (expected to be an ID)
	resolvedID, err := strconv.ParseInt(request.Identifier, 10, 64)
	if err != nil || resolvedID <= 0 {
		trustedService.Logger.Warnf("Invalid identifier format: %v", request.Identifier)
		handlers.WriteJSONResponse(w, handlers.ErrorResponse{Error: "Invalid identifier format"}, http.StatusBadRequest, trustedService.Logger)
		return
	}

	// Process removal based on trust status and entity type
	var removeErr error
	switch {
	case trustStatus == "trusted" && entityType == "character":
		removeErr = trustedService.RemoveTrustedCharacter(resolvedID)
	case trustStatus == "trusted" && entityType == "corporation":
		removeErr = trustedService.RemoveTrustedCorporation(resolvedID)
	case trustStatus == "untrusted" && entityType == "character":
		removeErr = trustedService.RemoveUntrustedCharacter(resolvedID)
	case trustStatus == "untrusted" && entityType == "corporation":
		removeErr = trustedService.RemoveUntrustedCorporation(resolvedID)
	default:
		trustedService.Logger.Warnf("Unsupported operation: trustStatus=%s, entityType=%s", trustStatus, entityType)
		handlers.WriteJSONResponse(w, handlers.ErrorResponse{Error: "Unsupported operation"}, http.StatusBadRequest, trustedService.Logger)
		return
	}

	if removeErr != nil {
		trustedService.Logger.Errorf("Error removing %s %s: %v", trustStatus, entityType, removeErr)
		handlers.WriteJSONResponse(w, handlers.ErrorResponse{Error: "Failed to remove entity"}, http.StatusInternalServerError, trustedService.Logger)
		return
	}

	trustedService.Logger.Infof("Validation succeeded: %s %s with ID %d successfully removed", trustStatus, entityType, resolvedID)
	handlers.WriteJSONResponse(w, handlers.SuccessResponse{Message: fmt.Sprintf("%s %s removed successfully", trustStatus, entityType)}, http.StatusOK, trustedService.Logger)
}

func RemoveTrustedCharacterHandler(s *handlers.SessionService, trustedService *service.TrustedService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleRemoveEntity(s, trustedService, w, r, "trusted", "character")
	}
}

func RemoveTrustedCorporationHandler(s *handlers.SessionService, trustedService *service.TrustedService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleRemoveEntity(s, trustedService, w, r, "trusted", "corporation")
	}
}

func RemoveUntrustedCharacterHandler(s *handlers.SessionService, trustedService *service.TrustedService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleRemoveEntity(s, trustedService, w, r, "untrusted", "character")
	}
}

func RemoveUntrustedCorporationHandler(s *handlers.SessionService, trustedService *service.TrustedService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleRemoveEntity(s, trustedService, w, r, "untrusted", "corporation")
	}
}

// Utility function to fetch entity data
func fetchEntityData(entityType string, data EntityData, token *oauth2.Token, esiService *service.EsiService) (EntityData, error) {
	if entityType == "character" {
		characterData, err := esiService.EsiClient.GetPublicCharacterData(data.ID, token)
		if err != nil {
			return EntityData{}, fmt.Errorf("error retrieving character data: %v", err)
		}
		corpID, err := esiService.EsiClient.GetCharacterCorporation(data.ID, token)
		if err != nil {
			return EntityData{}, fmt.Errorf("error retrieving character's corporation ID: %v", err)
		}
		corp, err := esiService.EsiClient.GetCorporationInfo(context.Background(), int(corpID))
		if err != nil {
			return EntityData{}, fmt.Errorf("error retrieving corporation info: %v", err)
		}
		data.Name = characterData.Name
		data.CorporationID = int64(corpID)
		data.CorporationName = corp.Name
		return data, nil
	} else if entityType == "corporation" {
		corp, err := esiService.EsiClient.GetCorporationInfo(context.Background(), int(data.ID))
		if err != nil {
			return EntityData{}, fmt.Errorf("error retrieving corporation name: %v", err)
		}
		data.Name = corp.Name
		return data, nil
	}
	return EntityData{}, fmt.Errorf("unknown entity type: %s", entityType)
}

func UpdateIsOnCouchHandler(s *handlers.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xlog.Logf("Update IsOnCouch Handler invoked")

		// Parse the request
		var request struct {
			ID        int64  `json:"id"`
			IsOnCouch bool   `json:"isOnCouch"`
			TableID   string `json:"tableId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			xlog.Logf("Error decoding JSON: %v", err)
			sendJSONError(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		xlog.Logf("Received Update IsOnCouch request: %+v", request)

		// Load the data
		data, err := persist.LoadTrustedCharacters()
		if err != nil {
			xlog.Logf("Error loading trusted characters: %v", err)
			sendJSONError(w, "Error loading trusted characters", http.StatusInternalServerError)
			return
		}

		// Update the IsOnCouch field based on the table and ID
		switch request.TableID {
		case "trusted-characters-table":
			for i, character := range data.TrustedCharacters {
				if character.CharacterID == request.ID {
					data.TrustedCharacters[i].IsOnCouch = request.IsOnCouch
					break
				}
			}
		case "trusted-corporations-table":
			for i, corporation := range data.TrustedCorporations {
				if corporation.CorporationID == request.ID {
					data.TrustedCorporations[i].IsOnCouch = request.IsOnCouch
					break
				}
			}
		default:
			xlog.Logf("Table ID was not recognized: %v", request.TableID)
			sendJSONError(w, "Error parsing tableID", http.StatusInternalServerError)
			return
		}

		// Save the updated data
		if err := persist.SaveTrustedCharacters(data); err != nil {
			xlog.Logf("Error saving trusted characters: %v", err)
			sendJSONError(w, "Error saving trusted characters", http.StatusInternalServerError)
			return
		}
		setDirtyBit(s, r, w)

		sendJSONResponse(w, http.StatusOK, map[string]string{"message": "IsOnCouch updated successfully"})

	}
}

func setDirtyBit(s *handlers.SessionService, r *http.Request, w http.ResponseWriter) {
	session, err := s.Get(r, handlers.SessionName)
	if err != nil {
		xlog.Logf("Failed to retrieve session for dirty bit update: %v", err)
		return
	}

	session.Values["trustedDataDirty"] = true
	if err := session.Save(r, w); err != nil {
		xlog.Logf("Failed to save session with updated dirty bit: %v", err)
	} else {
		xlog.Logf("Dirty bit set in session")
	}
}
