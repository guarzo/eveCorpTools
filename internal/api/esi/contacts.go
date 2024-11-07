package esi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/guarzo/zkillanalytics/internal/xlog"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strconv"
)

// AddContacts is a helper function to send contacts to the EVE API.
func (es *EsiClient) AddContacts(characterID int64, token *oauth2.Token, contactIDs []int64) error {
	// Prepare JSON payload
	contactIDsJSON, err := json.Marshal(contactIDs)
	if err != nil {
		xlog.Logf("Error encoding contact IDs: %v", err)
		return fmt.Errorf("error encoding contact IDs: %w", err)
	}

	// Build the request URL with query parameters
	baseURL := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/contacts/", characterID)
	params := url.Values{}
	params.Set("standing", strconv.FormatFloat(5.0, 'f', 1, 64))

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL+"?"+params.Encode(), bytes.NewBuffer(contactIDsJSON))
	if err != nil {
		xlog.Logf("Error creating request: %v", err)
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers for the request
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		xlog.Logf("Error executing request: %v", err)
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		xlog.Logf("Failed to add contacts: status %d", resp.StatusCode)
		return fmt.Errorf("failed to add contacts: received status %d", resp.StatusCode)
	}

	// Decode the response as an array of integers (contact IDs)
	var contacts []int
	if err := json.NewDecoder(resp.Body).Decode(&contacts); err != nil {
		xlog.Logf("Error decoding response body: %v", err)
		return fmt.Errorf("failed to decode response body: %v", err)
	}

	xlog.Logf("Contacts added successfully: %v", contacts)
	return nil
}

// DeleteContacts is a helper function to send contacts to the EVE API.
func (es *EsiClient) DeleteContacts(characterID int64, token *oauth2.Token, contactIDs []int64) error {
	// Prepare JSON payload
	contactIDsJSON, err := json.Marshal(contactIDs)
	if err != nil {
		xlog.Logf("Error encoding contact IDs: %v", err)
		return fmt.Errorf("error encoding contact IDs: %w", err)
	}

	// Build the request URL with query parameters
	baseURL := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/contacts/", characterID)
	params := url.Values{}
	for _, id := range contactIDs {
		params.Add("contact_ids", strconv.FormatInt(id, 10))
	}
	params.Set("datasource", "tranquility")

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", baseURL+"?"+params.Encode(), bytes.NewBuffer(contactIDsJSON))
	if err != nil {
		xlog.Logf("Error creating request: %v", err)
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers for the request
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		xlog.Logf("Error executing request: %v", err)
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusNoContent {
		xlog.Logf("Failed to delete contacts: status %d", resp.StatusCode)
		return fmt.Errorf("failed to delete contacts: received status %d", resp.StatusCode)
	}

	xlog.Logf("Contacts deleted successfully %v", contactIDs)
	return nil
}
