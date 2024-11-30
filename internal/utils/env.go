// internal/utils/config.todo

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
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

func GetVersion() string {
	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("No VERSION provided in environment, using placeholder")
		version = "v0.0.1"
	}
	return version
}

func GetESIEnv() (string, string, string) {
	clientID := os.Getenv("EVE_CLIENT_ID")
	clientSecret := os.Getenv("EVE_CLIENT_SECRET")
	callbackURL := os.Getenv("EVE_CALLBACK_URL")
	if clientID == "" || clientSecret == "" || callbackURL == "" {
		log.Fatalf("EVE_CLIENT_ID, EVE_CLIENT_SECRET, and EVE_CALLBACK_URL must be set")
	}
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
