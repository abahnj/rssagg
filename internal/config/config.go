package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Config represents the application configuration
type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Read loads the configuration from the config file
func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get config file path: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// SetUser updates the current user name in the config
func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	return write(*c)
}

// GetConfigFilePathFunc is the function type for getting config file path
type GetConfigFilePathFunc func() (string, error)

// GetConfigFilePath returns the full path to the config file
// It's exported to allow mocking in tests
var GetConfigFilePath = func() (string, error) {
	// For this project, we'll use the current directory instead of home
	currentDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	return filepath.Join(currentDir, configFileName), nil
}

// getConfigFilePath is an internal alias for GetConfigFilePath to maintain backward compatibility
var getConfigFilePath = GetConfigFilePath

// write saves the config to disk
func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

