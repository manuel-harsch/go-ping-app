// internal/config.go
package internal

import (
	"encoding/json"
	"io"
	"os"
)

// Config represents the structure of the JSON configuration file
type Config struct {
	Host        string `json:"host"`
	CycleTime   int    `json:"cycle_time_seconds"` // Ping cycle time in seconds
	PingTimeout int    `json:"ping_timeout"`       // Ping timeout in milliseconds
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to a JSON file
func SaveConfig(config *Config, filePath string) error {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Host:        "8.8.8.8", // Default to Google Public DNS
		CycleTime:   5,         // 5 seconds between pings
		PingTimeout: 1000,      // 1000 ms (1 second) ping timeout
	}
}
