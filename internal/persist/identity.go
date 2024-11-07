package persist

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

// Directory where identity files are stored
const trustDirectory = "data/trust"

// GetMainIdentityToken retrieves the token for the main identity.
func GetMainIdentityToken(mainIdentity int64) (oauth2.Token, error) {
	identities, err := LoadIdentities(mainIdentity)
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("unable to retrieve token for main identity")
	}
	token, exists := identities.Tokens[fmt.Sprintf("%d", mainIdentity)]
	if !exists {
		return oauth2.Token{}, fmt.Errorf("main identity token not found")
	}
	return token, nil
}

// LoadIdentityToken retrieves the token for a specified character.
func LoadIdentityToken(mainIdentity, characterID int64) (oauth2.Token, error) {
	identities, err := LoadIdentities(mainIdentity)
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("unable to retrieve token for character %d", characterID)
	}

	token, exists := identities.Tokens[fmt.Sprintf("%d", characterID)]
	if !exists {
		return oauth2.Token{}, fmt.Errorf("token not found for character %d", characterID)
	}
	return token, nil
}

func LoadIdentities(mainIdentity int64) (*model.Identities, error) {
	if mainIdentity == 0 {
		return nil, fmt.Errorf("logged in user not provided")
	}

	xlog.Logf("Loading identities for mainIdentity: %d", mainIdentity)

	mainIdentityStr := fmt.Sprintf("%d", mainIdentity)

	identities := &model.Identities{
		MainIdentity: mainIdentityStr,
		Tokens:       make(map[string]oauth2.Token),
	}

	identityFile := getIdentityFileName(mainIdentity)
	if _, err := os.Stat(identityFile); os.IsNotExist(err) {
		xlog.Log("No identity file found. Initializing new Identities.")
		return identities, nil
	}

	// Attempt to load with the new model (map[string]oauth2.Token)
	if err := DecryptData(identityFile, identities); err == nil {
		xlog.Logf("Loaded identities with new model: %s", SafeLogIdentities(identities))
		return identities, nil
	}

	// Fallback: Attempt to load with the old model (map[int64]oauth2.Token)
	xlog.Log("Attempting to load identities with old model format.")
	oldIdentities := struct {
		MainIdentity string                 `json:"main_identity"`
		Tokens       map[int64]oauth2.Token `json:"identities"`
	}{}

	if err := DecryptData(identityFile, &oldIdentities); err != nil {
		xlog.Logf("Error decrypting identities file with old model: %v", err)
		return nil, err
	}

	// Convert old model (map[int64]) to new model (map[string])
	for k, v := range oldIdentities.Tokens {
		identities.Tokens[fmt.Sprintf("%d", k)] = v
	}

	// Save converted identities with the new model format
	if saveErr := SaveIdentities(mainIdentity, identities); saveErr != nil {
		xlog.Logf("Error saving identities after conversion: %v", saveErr)
		return nil, saveErr
	}

	xlog.Logf("Converted and saved identities to new model: %s", SafeLogIdentities(identities))
	return identities, nil
}

func SaveIdentities(mainIdentity int64, ids *model.Identities) error {
	if mainIdentity == 0 {
		return fmt.Errorf("no main identity provided")
	}
	xlog.Logf("Saving identities for mainIdentity: %d", mainIdentity)

	xlog.Logf("Saving identities: %s", SafeLogIdentities(ids))
	return EncryptData(getIdentityFileName(mainIdentity), ids)
}

// getIdentityFileName generates the file path for a given main identity.
func getIdentityFileName(mainIdentity int64) string {
	return filepath.Join(trustDirectory, fmt.Sprintf("%d_identity.json", mainIdentity))
}

// UpdateIdentities loads, updates, and saves identities.
func UpdateIdentities(mainIdentity int64, updateFunc func(*model.Identities) error) error {
	xlog.Logf("Loading identities for mainIdentity: %d", mainIdentity)
	ids, err := LoadIdentities(mainIdentity)
	if err != nil {
		xlog.Logf("Error in LoadIdentities: %v", err)
		return err
	}

	// Log current identities without tokens
	xlog.Logf("Current identities: %s", SafeLogIdentities(ids))

	if err = updateFunc(ids); err != nil {
		xlog.Logf("Error in updateFunc: %v", err)
		return err
	}

	// Log identities after update without tokens
	xlog.Logf("Identities after updateFunc: %s", SafeLogIdentities(ids))

	if err = SaveIdentities(mainIdentity, ids); err != nil {
		xlog.Logf("Error in SaveIdentities: %v", err)
		return err
	}

	xlog.Log("Identities successfully updated and saved.")
	return nil
}

// SafeLogIdentities returns a string representation of Identities without sensitive tokens.
func SafeLogIdentities(ids *model.Identities) string {
	var tokensInfo []string
	for charID, token := range ids.Tokens {
		tokensInfo = append(tokensInfo, fmt.Sprintf("CharacterID: %s, TokenType: %s, Expiry: %s", charID, token.TokenType, token.Expiry))
	}
	return fmt.Sprintf("MainIdentity: %s, Tokens: [%s]", ids.MainIdentity, strings.Join(tokensInfo, "; "))
}

// DeleteIdentity deletes the identity file for a specified main identity.
func DeleteIdentity(mainIdentity int64) error {
	return os.Remove(getIdentityFileName(mainIdentity))
}
