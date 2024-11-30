// internal/utils/config.todo

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

// GetPort retrieves the port from the PORT environment variable.
func GetPort() int {
	portStr, exists := os.LookupEnv("PORT")
	if !exists || portStr == "" {
		portStr = "8080"
		log.Printf("PORT not set. Using default port %s", portStr)
	} else {
		log.Printf("Using port from environment: %s", portStr)
	}

	var port int
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}
	return port
}

// GetUserAgent retrieves the USER_AGENT from the environment variable.
func GetUserAgent() string {
	userAgent, exists := os.LookupEnv("USER_AGENT")
	if !exists {
		userAgent = "placeholder@gmail.com"
		log.Printf("USER_AGENT not set. Using default: %s", userAgent)

	} else {
		log.Printf("Using USER_AGENT from environment: %s", userAgent)
	}
	return userAgent
}

func GetHostConfig() string {
	hostConfig := os.Getenv("HOST_CONFIG")
	if hostConfig == "" {
		fmt.Println("No HOST_CONFIG provided in environment, using placeholder")
		hostConfig = "tps.zoolanders.space"
	}
	return hostConfig
}

func GetHost(host string) string {
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return host
}

func GetVersion() string {
	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("No VERSION provided in environment, using placeholder")
		version = "v0.0.1"
	}
	return version
}

func GetESIEnv(host string) (string, string, string) {
	// Define the unique environment variable names for the specific host
	clientID := os.Getenv(fmt.Sprintf("%s_EVE_CLIENT_ID", host))
	clientSecret := os.Getenv(fmt.Sprintf("%s_EVE_CLIENT_SECRET", host))
	callbackURL := os.Getenv(fmt.Sprintf("%s_EVE_CALLBACK_URL", host))

	// Ensure that all required environment variables are set
	if clientID == "" || clientSecret == "" || callbackURL == "" {
		log.Printf("clientid :%s, clientsecret :%s, callbackurl %s", clientID, clientSecret, callbackURL)
		log.Fatalf("%s_EVE_CLIENT_ID, %s_EVE_CLIENT_SECRET, and %s_EVE_CALLBACK_URL must be set", host, host, host)
	}

	// Return the environment variables
	return clientID, clientSecret, callbackURL
}

func GetSecretKey() (string, []byte) {
	secret := os.Getenv("SECRET_KEY")
	var key []byte
	var err error
	if secret == "" {
		key, err = generateSecret()
		if err != nil {
			log.Fatalf("Failed to generate key: %v", err)
		}
		secret = base64.StdEncoding.EncodeToString(key)
		log.Printf("Generated key: %s -- this should only be used for testing", secret)
	} else {
		key, err = base64.StdEncoding.DecodeString(secret)
		if err != nil {
			log.Fatalf("Failed to decode key: %v", err)
		}
	}
	return secret, key
}

func generateSecret() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
