// internal/utils/config.go

package utils

import (
	"fmt"
	"log"
	"os"
)

// GetPort retrieves the port from the PORT environment variable.
// If not set or empty, it defaults to the provided defaultPort.
func GetPort(defaultPort string) int {
	portStr, exists := os.LookupEnv("PORT")
	if !exists || portStr == "" {
		log.Printf("PORT not set. Using default port %s", defaultPort)
		portStr = defaultPort
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
// If not set, it defaults to the provided defaultUA.
func GetUserAgent(defaultUA string) string {
	userAgent, exists := os.LookupEnv("USER_AGENT")
	if !exists {
		log.Printf("USER_AGENT not set. Using default: %s", defaultUA)
		userAgent = defaultUA
	} else {
		log.Printf("Using USER_AGENT from environment: %s", userAgent)
	}
	return userAgent
}
