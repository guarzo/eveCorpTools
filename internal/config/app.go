package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/guarzo/zkillanalytics/internal/utils"
)

// AppSetup holds all environment variables for easy access and passing
type AppSetup struct {
	Port         int
	UserAgent    string
	HostConfig   string
	Version      string
	ClientID     string
	ClientSecret string
	CallbackURL  string
	Key          []byte
	Secret       string
}

// NewAppSetup initializes and returns a Config struct with values from environment variables
func NewAppSetup() (*AppSetup, error) {

	// Step 1: Check if the PORT environment variable is set before loading the .env file
	originalPort := os.Getenv("PORT")

	// Step 2: Load the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// Step 3: If PORT is not set by the .env file, restore it to the original value
	if originalPort != "" && os.Getenv("PORT") == "" {
		// Restore the original PORT value if it was not set by the .env file
		err := os.Setenv("PORT", originalPort)
		if err != nil {
			fmt.Println("Error setting original port")
		}
	}

	port := utils.GetPort()
	userAgent := utils.GetUserAgent()

	hostConfig := utils.GetHostConfig()

	version := utils.GetVersion()

	clientID, clientSecret, callbackURL := utils.GetESIEnv()

	// Handle SECRET_KEY, generate if missing
	secret, key := utils.GetSecretKey()

	return &AppSetup{
		Port:         port,
		UserAgent:    userAgent,
		HostConfig:   hostConfig,
		Version:      version,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		CallbackURL:  callbackURL,
		Key:          key,
		Secret:       secret,
	}, nil
}
