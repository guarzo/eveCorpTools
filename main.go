// main.go

package main

import (
	"log"

	"github.com/guarzo/zkillanalytics/cmd"
	"github.com/guarzo/zkillanalytics/internal/config"
)

func main() {
	appSetup, err := config.NewAppSetup()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	// Start the web server
	cmd.StartServer(appSetup)
}
