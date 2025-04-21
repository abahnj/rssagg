package cli

import (
	"errors"
	"testing"
)

func TestCommands_Register(t *testing.T) {
	cmds := NewCommands()
	
	// Register a mock handler
	mockHandler := func(*State, Command) error { return nil }
	cmds.Register("test", mockHandler)
	
	// Verify the handler was registered
	if _, exists := cmds.handlers["test"]; !exists {
		t.Error("Expected handler to be registered")
	}
}

func TestCommands_Run(t *testing.T) {
	cmds := NewCommands()
	state := &State{}
	
	t.Run("Unknown command", func(t *testing.T) {
		cmd := Command{Name: "unknown"}
		err := cmds.Run(state, cmd)
		
		if err == nil {
			t.Error("Expected error for unknown command")
		}
	})
	
	t.Run("Run existing command", func(t *testing.T) {
		// Mock handler that returns nil
		calledHandler := false
		mockHandler := func(*State, Command) error {
			calledHandler = true
			return nil
		}
		
		cmds.Register("test", mockHandler)
		cmd := Command{Name: "test"}
		err := cmds.Run(state, cmd)
		
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if !calledHandler {
			t.Error("Expected handler to be called")
		}
	})
	
	t.Run("Handle error from command", func(t *testing.T) {
		expectedErr := errors.New("test error")
		mockHandler := func(*State, Command) error {
			return expectedErr
		}
		
		cmds.Register("error", mockHandler)
		cmd := Command{Name: "error"}
		err := cmds.Run(state, cmd)
		
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}