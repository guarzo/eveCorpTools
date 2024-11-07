package config

import (
	"fmt"
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
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
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
