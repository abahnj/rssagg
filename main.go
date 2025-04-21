package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abahnj/rssagg/internal/cli"
	"github.com/abahnj/rssagg/internal/config"
)

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Create application state
	state := &cli.State{
		Config: &cfg,
	}

	// Set up commands
	commands := cli.NewCommands()
	registerCommands(commands)

	// Process command line arguments
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: command name is required")
		fmt.Println("Usage: rssagg <command> [arguments]")
		os.Exit(1)
	}

	// Parse command
	cmdName := cli.CommandName(args[1])
	cmdArgs := args[2:]
	cmd := cli.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	// Run the command
	if err := commands.Run(state, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}