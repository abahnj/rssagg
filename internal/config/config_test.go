package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the getConfigFilePath function for testing
	originalGetConfigFilePath := getConfigFilePath
	defer func() {
		getConfigFilePath = originalGetConfigFilePath
	}()

	testConfigPath := filepath.Join(tempDir, configFileName)
	getConfigFilePath = func() (string, error) {
		return testConfigPath, nil
	}

	t.Run("Read nonexistent file", func(t *testing.T) {
		_, err := Read()
		if err == nil {
			t.Fatal("Expected error when reading nonexistent file, got nil")
		}
	})

	t.Run("Read invalid JSON", func(t *testing.T) {
		// Write invalid JSON to the config file
		if err := os.WriteFile(testConfigPath, []byte("invalid json"), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		_, err := Read()
		if err == nil {
			t.Fatal("Expected error when reading invalid JSON, got nil")
		}
	})

	t.Run("Read valid config", func(t *testing.T) {
		// Write valid config to the file
		expected := Config{
			DBURL:          "postgres://test",
			CurrentUserName: "testuser",
		}
		data, err := json.MarshalIndent(expected, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal test config: %v", err)
		}

		if err := os.WriteFile(testConfigPath, data, 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		cfg, err := Read()
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}

		if cfg.DBURL != expected.DBURL {
			t.Errorf("Expected DBURL=%q, got %q", expected.DBURL, cfg.DBURL)
		}

		if cfg.CurrentUserName != expected.CurrentUserName {
			t.Errorf("Expected CurrentUserName=%q, got %q", expected.CurrentUserName, cfg.CurrentUserName)
		}
	})

	t.Run("SetUser", func(t *testing.T) {
		// Start with a clean config
		initialCfg := Config{
			DBURL: "postgres://test",
		}
		data, err := json.MarshalIndent(initialCfg, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal test config: %v", err)
		}

		if err := os.WriteFile(testConfigPath, data, 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		// Read the config
		cfg, err := Read()
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}

		// Set the user
		newUsername := "newuser"
		if err := cfg.SetUser(newUsername); err != nil {
			t.Fatalf("Failed to set user: %v", err)
		}

		// Read the config again to verify changes
		updatedCfg, err := Read()
		if err != nil {
			t.Fatalf("Failed to read updated config: %v", err)
		}

		if updatedCfg.CurrentUserName != newUsername {
			t.Errorf("Expected CurrentUserName=%q after SetUser, got %q", newUsername, updatedCfg.CurrentUserName)
		}

		if updatedCfg.DBURL != initialCfg.DBURL {
			t.Errorf("Expected DBURL to remain %q, got %q", initialCfg.DBURL, updatedCfg.DBURL)
		}
	})

	t.Run("JSON Marshal/Unmarshal", func(t *testing.T) {
		cfg := Config{
			DBURL:          "postgres://test",
			CurrentUserName: "testuser",
		}

		data, err := json.Marshal(cfg)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		var unmarshaled Config
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal config: %v", err)
		}

		if unmarshaled.DBURL != cfg.DBURL {
			t.Errorf("Expected DBURL=%q after marshal/unmarshal, got %q", cfg.DBURL, unmarshaled.DBURL)
		}

		if unmarshaled.CurrentUserName != cfg.CurrentUserName {
			t.Errorf("Expected CurrentUserName=%q after marshal/unmarshal, got %q", cfg.CurrentUserName, unmarshaled.CurrentUserName)
		}

		// Verify JSON field names
		expectedJSON := `{"db_url":"postgres://test","current_user_name":"testuser"}`
		if string(data) != expectedJSON {
			t.Errorf("Expected JSON=%q, got %q", expectedJSON, string(data))
		}
	})
}

func TestGetConfigFilePath(t *testing.T) {
	// Since we can't mock os.Getwd directly, we'll test the actual function behavior
	
	// The actual getConfigFilePath function should:
	// 1. Get the current directory
	// 2. Join it with the config file name
	
	// We can test that it at least returns a non-empty path without error
	path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("getConfigFilePath failed: %v", err)
	}
	
	if path == "" {
		t.Error("Expected non-empty path, got empty string")
	}
	
	// Check that the path ends with the config file name
	if filepath.Base(path) != configFileName {
		t.Errorf("Expected path to end with %q, got %q", configFileName, filepath.Base(path))
	}
}