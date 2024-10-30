// cmd/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
		// If the config file is corrupted, we should exit the program with a non-zero exit code
		os.Exit(1)
	} else {
		log.Printf("Loaded Config: Host=%s, CycleTime=%d, PingTimeout=%d", cfg.Host, cfg.CycleTime, cfg.PingTimeout)
	}

	// Set up the Gin web server
	router := gin.Default()

	// Serve static files (the GUI)
	router.Static("/static", "./assets/static")

	// API: Fetch current configuration
	router.GET("/api/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, cfg)
	})

	// API: Update the configuration
	router.POST("/api/config", func(c *gin.Context) {
		var newConfig internal.Config
		if err := c.BindJSON(&newConfig); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration data"})
			return
		}

		// Save the new configuration to the file
		if err := internal.SaveConfig(&newConfig, configFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
			return
		}

		// Update the in-memory config (so changes apply without restart)
		cfg = &newConfig
		c.JSON(http.StatusOK, cfg)
	})

	// API: Start ping process
	router.POST("/api/ping/start", func(c *gin.Context) {
		// Logic to start the ping process would go here
		c.JSON(http.StatusOK, gin.H{"message": "Ping started"})
	})

	// API: Stop ping process
	router.POST("/api/ping/stop", func(c *gin.Context) {
		// Logic to stop the ping process would go here
		c.JSON(http.StatusOK, gin.H{"message": "Ping stopped"})
	})

	// Start the web server
	router.Run(":8080")
}
