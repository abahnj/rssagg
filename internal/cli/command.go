package cli

import (
	"fmt"
)

// CommandName is a custom string type for command names
type CommandName string

// Command represents a CLI command with its arguments
type Command struct {
	Name CommandName
	Args []string
}

// HandlerFunc is the function signature for command handlers
type HandlerFunc func(*State, Command) error

// Commands manages the CLI commands and their handlers
type Commands struct {
	handlers map[CommandName]HandlerFunc
}

// NewCommands creates a new Commands instance
func NewCommands() *Commands {
	return &Commands{
		handlers: make(map[CommandName]HandlerFunc),
	}
}

// Register adds a new command handler
func (c *Commands) Register(name CommandName, handler HandlerFunc) {
	c.handlers[name] = handler
}

// Run executes the appropriate handler for the given command
func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(s, cmd)
}