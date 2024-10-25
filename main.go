package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gambtho/zkillanalytics/cmd"
)

func main() {
	// Read the port from the environment variable
	portStr := os.Getenv("PORT")

	// Default to 8081 if the environment variable is not set
	port := 8081
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Printf("Invalid port number: %s\n", portStr)
			os.Exit(1)
		}
		port = p
	}

	// Start the web server
	cmd.StartServer(port)
}
