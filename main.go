// main.go

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/guarzo/zkillanalytics/cmd"
)

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get the port number from the environment variable
	portStr := os.Getenv("PORT")

	// Default the port to 8081 if the environment variable is not set
	port := 8081
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err != nil {
			fmt.Printf("Invalid port number: %s\n", portStr)
			os.Exit(1)
		} else {
			port = p
		}
	} else {
		fmt.Println("PORT environment variable not set. Defaulting to 8081.")
	}

	userAgent := os.Getenv("USER_AGENT")
	if userAgent == "" {
		fmt.Println("No userAgent provided in environment, using placeholder")
		userAgent = "placeholder@gmail.com"
	}

	hostConfig := os.Getenv("HOST_CONFIG")
	if hostConfig == "" {
		fmt.Println("No hostConfig provided in environment, using placeholder")
		userAgent = "tps.zoolanders.space"
	}

	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("No version provided in environment, using placeholder")
		userAgent = "placeholder@gmail.com"
	}

	// Start the web server
	cmd.StartServer(port, userAgent, version, hostConfig)
}
