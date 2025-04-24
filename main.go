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
		fmt.Println("Available commands:")
		fmt.Println("  login <username>  - Log in as a user")
		fmt.Println("  register <username> - Register a new user")
		fmt.Println("  users - List all users")
		fmt.Println("  reset - Delete all users")
		fmt.Println("  feeds - List all feeds")
		fmt.Println("  addfeed <name> <url> - Add a new feed")
		fmt.Println("  follow <url> - Follow an existing feed")
		fmt.Println("  unfollow <url> - Unfollow a feed")
		fmt.Println("  following - List feeds you're following")
		fmt.Println("  browse [limit] - View posts from feeds you follow (default limit: 10)")
		fmt.Println("  agg <duration> - Aggregate and show feed content every <duration> (e.g. 30s, 1m)")
		os.Exit(0)
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