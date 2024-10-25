package routes

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

func getAPIKey() (string, error) {
	ApiKey := os.Getenv("API_KEY")
	if ApiKey == "" {
		log.Printf("API key not found in environment variable, attempting to read from file") // Log the attempt
		const apiKeyPath = "./apikey.txt"                                                     // Ensure this path is correct

		absPath, err := os.Getwd()
		if err != nil {
			log.Printf("Error getting current working directory: %v", err)
			return "", err
		}
		fullPath := absPath + "/" + apiKeyPath

		log.Printf("Attempting to read API key from: %s", fullPath) // Log the full path

		key, err := os.ReadFile(fullPath)
		if err != nil {
			log.Printf("Error reading API key from %s: %v", fullPath, err)
			return "", err
		}
		ApiKey = strings.TrimSpace(string(key))
	}

	return ApiKey, nil
}

func AppraiseLootHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	loot := string(body)
	// log.Printf("Received loot data: %s", loot) // Debugging

	apiKey, err := getAPIKey()
	if err != nil {
		log.Printf("Unable to read API key: %v", err)
		http.Error(w, "Unable to read API key", http.StatusInternalServerError)
		return
	}

	url := "https://janice.e-351.com/api/rest/v2/appraisal?market=2&designation=appraisal&pricing=buy&pricingVariant=immediate&persist=true&compactize=true&pricePercentage=1"
	req, err := http.NewRequest("POST", url, strings.NewReader(loot))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("X-ApiKey", apiKey)

	//log.Printf("Request URL: %s", url)            // Log the URL
	//log.Printf("Request Headers: %v", req.Header) // Log request headers
	//log.Printf("Request Body: %s", loot)          // Log request body

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	//log.Printf("Response status from external API: %s", resp.Status) // Debugging

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log.Printf("Response body from external API: %s", string(respBody)) // Debugging

	var result struct {
		ImmediatePrices struct {
			TotalBuyPrice float64 `json:"totalBuyPrice"`
		} `json:"immediatePrices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		response := map[string]float64{"totalBuyPrice": 0}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Round down the total buy price
	roundedTotalBuyPrice := math.Floor(result.ImmediatePrices.TotalBuyPrice)

	// log.Printf("Total Buy Price extracted: %f", result.ImmediatePrices.TotalBuyPrice) // Debugging

	response := map[string]float64{"totalBuyPrice": roundedTotalBuyPrice}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
