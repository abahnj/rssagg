package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abahnj/rssagg/internal/cli"
	"github.com/abahnj/rssagg/internal/config"
	"github.com/abahnj/rssagg/internal/database"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	conn, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Create application state
	state := &cli.State{
		Config: &cfg,
		Db:     database.New(conn),
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
