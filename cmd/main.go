// cmd/main.go
package main

import (
	"log"
	"os"

	"github.com/manuel-harsch/go-ping-app/internal"
)

const configFilePath = "config.json"

func main() {
	// Check if the configuration file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// File does not exist, create a default configuration
		log.Println("Config file not found, creating default config.")
		defaultConfig := internal.DefaultConfig() // Call DefaultConfig from internal package
		if err := internal.SaveConfig(defaultConfig, configFilePath); err != nil {
			log.Fatalf("Failed to create default config: %v", err)
		}
	}

	// Load the configuration
	cfg, err := internal.LoadConfig(configFilePath) // Call LoadConfig from internal package
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Printf("Loaded Config: Host=%s, CycleTime=%d, PingTimeout=%d", cfg.Host, cfg.CycleTime, cfg.PingTimeout)

	// Now you can use the loaded configuration for further processing
}
