package main

import (
	"errors"
	"fmt"

	"github.com/abahnj/rssagg/internal/cli"
)

// ErrMissingUsername is returned when the login command doesn't have a username argument
var ErrMissingUsername = errors.New("username is required for login")

// registerCommands sets up all available commands
func registerCommands(commands *cli.Commands) {
	commands.Register("login", handlerLogin)
}

// handlerLogin handles the login command
func handlerLogin(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) < 1 {
		return ErrMissingUsername
	}

	// Check for nil config
	if s.Config == nil {
		return errors.New("config is not initialized")
	}

	username := cmd.Args[0]
	if err := s.Config.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("User set to: %s\n", username)
	return nil
}