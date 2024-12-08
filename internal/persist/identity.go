package persist

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"

	"github.com/guarzo/zkillanalytics/internal/model"
	"github.com/guarzo/zkillanalytics/internal/xlog"
)

// GetMainIdentityToken retrieves the token for the main identity.
func GetMainIdentityToken(mainIdentity int64, host string) (oauth2.Token, error) {
	identities, err := LoadIdentities(mainIdentity, host)
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
func LoadIdentityToken(mainIdentity, characterID int64, host string) (oauth2.Token, error) {
	identities, err := LoadIdentities(mainIdentity, host)
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("unable to retrieve token for character %d", characterID)
	}

	token, exists := identities.Tokens[fmt.Sprintf("%d", characterID)]
	if !exists {
		return oauth2.Token{}, fmt.Errorf("token not found for character %d", characterID)
	}
	return token, nil
}

func LoadIdentities(mainIdentity int64, host string) (*model.Identities, error) {
	if mainIdentity == 0 {
		return nil, fmt.Errorf("logged in user not provided")
	}

	xlog.Logf("Loading identities for mainIdentity: %d from host: %s", mainIdentity, host)

	// Use the host to determine the directory
	identityFile := getIdentityFileName(mainIdentity, host)
	if _, err := os.Stat(identityFile); os.IsNotExist(err) {
		xlog.Log("No identity file found. Initializing new Identities.")
		return &model.Identities{
			MainIdentity: fmt.Sprintf("%d", mainIdentity),
			Tokens:       make(map[string]oauth2.Token),
		}, nil
	}

	// Attempt to load with the new model
	var identities model.Identities
	if err := DecryptData(identityFile, &identities); err == nil {
		xlog.Logf("Loaded identities: %s", SafeLogIdentities(&identities))
		return &identities, nil
	}
	return nil, errors.New("unable to decrypt identities")
}

func SaveIdentities(mainIdentity int64, ids *model.Identities, host string) error {
	if mainIdentity == 0 {
		return fmt.Errorf("no main identity provided")
	}
	xlog.Logf("Saving identities for mainIdentity: %d to host: %s", mainIdentity, host)

	// Save the identities to the correct directory based on the host
	identityFile := getIdentityFileName(mainIdentity, host)
	xlog.Logf("Saving identities to: %s", identityFile)

	// Encrypt and save data
	return EncryptData(identityFile, ids)
}

// getIdentityFileName generates the file path for a given main identity, based on the host.
func getIdentityFileName(mainIdentity int64, host string) string {
	var subAppDirectory string

	switch host {
	case "loot.zoolanders.space":
		subAppDirectory = "data/loot"
	case "tps.zoolanders.space":
		subAppDirectory = "data/tps"
	case "trust.zoolanders.space":
		subAppDirectory = "data/trust"
	case "localhost":
		subAppDirectory = "data/trust"
	default:
		subAppDirectory = "data/default" // You could add a default case or handle this as an error
	}

	return filepath.Join(subAppDirectory, fmt.Sprintf("%d_identity.json", mainIdentity))
}

// UpdateIdentities loads, updates, and saves identities.
func UpdateIdentities(mainIdentity int64, host string, updateFunc func(*model.Identities) error) error {
	xlog.Logf("Loading identities for mainIdentity: %d, host: %s", mainIdentity, host)
	ids, err := LoadIdentities(mainIdentity, host)
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

	if err = SaveIdentities(mainIdentity, ids, host); err != nil {
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
func DeleteIdentity(mainIdentity int64, host string) error {
	return os.Remove(getIdentityFileName(mainIdentity, host))
}
