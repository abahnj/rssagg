package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/abahnj/rssagg/internal/cli"
	"github.com/abahnj/rssagg/internal/config"
)

func TestHandlerLogin(t *testing.T) {
	// For tests that don't need a real config
	t.Run("Missing username", func(t *testing.T) {
		cmd := cli.Command{
			Name: "login",
			Args: []string{},
		}
		// For this test, we don't need a real config since we'll hit the error before using it
		state := &cli.State{}
		
		err := handlerLogin(state, cmd)
		
		if !errors.Is(err, ErrMissingUsername) {
			t.Errorf("Expected ErrMissingUsername, got %v", err)
		}
	})
	
	// Test with nil config to trigger an error
	t.Run("Set user error with nil config", func(t *testing.T) {
		cmd := cli.Command{
			Name: "login",
			Args: []string{"erroruser"},
		}
		state := &cli.State{
			Config: nil, // Nil config should trigger an error
		}
		
		err := handlerLogin(state, cmd)
		
		if err == nil {
			t.Error("Expected error with nil config, got nil")
		}
		
		if err != nil && err.Error() != "config is not initialized" {
			t.Errorf("Expected 'config is not initialized' error, got: %v", err)
		}
	})
	
	// For tests that need a real config, use a temporary file
	t.Run("Integration tests with real config", func(t *testing.T) {
		// Create a temp directory for tests
		tempDir, err := os.MkdirTemp("", "login-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)
		
		// Create a temporary config file
		tempConfigPath := filepath.Join(tempDir, ".gatorconfig.json")
		initialConfig := `{"db_url": "postgres://test"}`
		if err := os.WriteFile(tempConfigPath, []byte(initialConfig), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}
		
		// Save original config path function and restore it after test
		originalGetConfigPath := config.GetConfigFilePath
		defer func() {
			config.GetConfigFilePath = originalGetConfigPath
		}()
		
		// Override config path for test
		config.GetConfigFilePath = func() (string, error) {
			return tempConfigPath, nil
		}
		
		// Test setting user successfully
		t.Run("Set user successfully", func(t *testing.T) {
			// Read the test config
			cfg, err := config.Read()
			if err != nil {
				t.Fatalf("Failed to read test config: %v", err)
			}
			
			// Setup command and state
			cmd := cli.Command{
				Name: "login",
				Args: []string{"testuser"},
			}
			state := &cli.State{
				Config: &cfg,
			}
			
			// Execute handler
			err = handlerLogin(state, cmd)
			
			// Verify results
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			// Read config again to verify changes
			updatedCfg, err := config.Read()
			if err != nil {
				t.Fatalf("Failed to read updated config: %v", err)
			}
			
			if updatedCfg.CurrentUserName != "testuser" {
				t.Errorf("Expected username to be 'testuser', got %q", updatedCfg.CurrentUserName)
			}
		})
	})
}