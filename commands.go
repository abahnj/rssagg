package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/abahnj/rssagg/internal/cli"
	"github.com/abahnj/rssagg/internal/database"
	"github.com/google/uuid"
)

// ErrMissingUsername is returned when the login command doesn't have a username argument
var ErrMissingUsername = errors.New("username is required for login")

// RSS feed structures
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// registerCommands sets up all available commands
func registerCommands(commands *cli.Commands) {
	commands.Register("login", handlerLogin)
	commands.Register("register", handlerRegister)
	commands.Register("reset", handlerDeleteAllUsers)
	commands.Register("users", handlerListUsers)
	commands.Register("agg", handlerAggregator)
	commands.Register("addfeed", handlerAddFeed)
	commands.Register("feeds", handlerListFeeds)
	commands.Register("follow", handlerFollowFeed)
	commands.Register("following", handlerListFollowing)
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
	ctx := context.Background()

	_, err := s.Db.GetUserByName(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to fetch  user with name %s: %w", username, err)
	}

	if err := s.Config.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("User  %s logged in\n", username)
	return nil
}

func handlerRegister(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("username is required for register")
	}

	// Check for nil config
	if s.Config == nil {
		return errors.New("config is not initialized")
	}

	username := cmd.Args[0]
	ctx := context.Background()
	createUserParams := database.CreateUserParams{
		ID:   uuid.New(),
		Name: username,
	}

	_, err := s.Db.CreateUser(ctx, createUserParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return handlerLogin(s, cmd)

}

func handlerDeleteAllUsers(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	if err := s.Db.DeleteAllUsers(ctx); err != nil {
		return fmt.Errorf("failed to delete all users: %w", err)
	}
	fmt.Println("All users deleted")
	return nil
}

func handlerListUsers(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	
	currentUserName := ""
	if s.Config != nil {
		currentUserName = s.Config.CurrentUserName
	}
	
	for _, user := range users {
		if user.Name == currentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	
	return nil
}

// fetchFeed retrieves and parses an RSS feed from the given URL
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set User-Agent header to identify our program
	req.Header.Set("User-Agent", "gator")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-success status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the XML
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("error parsing XML: %w", err)
	}

	// Unescape HTML entities in the channel's title and description
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Unescape HTML entities in each item's title and description
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

// handlerAggregator handles the agg command
func handlerAggregator(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	
	// URL of the feed to fetch
	feedURL := "https://www.wagslane.dev/index.xml"
	
	// Fetch the feed
	feed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	
	// Print the feed details
	fmt.Printf("Feed: %s\n", feed.Channel.Title)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Printf("Items: %d\n\n", len(feed.Channel.Item))
	
	// Print each item
	for i, item := range feed.Channel.Item {
		fmt.Printf("--- Item %d ---\n", i+1)
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Published: %s\n", item.PubDate)
		fmt.Printf("Description: %s\n\n", item.Description)
	}
	
	return nil
}

// handlerAddFeed handles the addfeed command to add a new RSS feed
func handlerAddFeed(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) < 2 {
		return errors.New("both feed name and URL are required")
	}

	// Get the current user
	if s.Config == nil || s.Config.CurrentUserName == "" {
		return errors.New("you must be logged in to add a feed")
	}

	ctx := context.Background()
	
	// Get the current user from the database
	user, err := s.Db.GetUserByName(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	
	// Check if feed already exists
	existingFeed, err := s.Db.GetFeedByURL(ctx, feedURL)
	var feed database.Feed
	
	if err == nil {
		// Feed already exists, we'll use this one
		feed = existingFeed
		fmt.Printf("Feed already exists, using existing feed: %s\n", feed.Name)
	} else {
		// Create a new feed record
		createFeedParams := database.CreateFeedParams{
			ID:     uuid.New(),
			Name:   feedName,
			Url:    feedURL,
			UserID: user.ID,
		}
		
		feed, err = s.Db.CreateFeed(ctx, createFeedParams)
		if err != nil {
			return fmt.Errorf("failed to create feed: %w", err)
		}
		
		// Print out the new feed details
		fmt.Printf("Feed added successfully:\n")
		fmt.Printf("  ID: %s\n", feed.ID)
		fmt.Printf("  Name: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  Created: %s\n", feed.CreatedAt.Time.Format("2006-01-02 15:04:05"))
	}
	
	// Now create a feed follow for the current user
	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	}
	
	feedFollow, err := s.Db.CreateFeedFollow(ctx, createFeedFollowParams)
	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}
	
	fmt.Printf("\nYou are now following this feed as user %s\n", feedFollow.UserName)
	
	return nil
}

// handlerListFeeds handles the feeds command to list all feeds
func handlerListFeeds(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	
	// Fetch all feeds with user information
	feeds, err := s.Db.GetFeedsWithUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
	}
	
	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}
	
	fmt.Printf("Found %d feeds:\n\n", len(feeds))
	
	// Display each feed with user information
	for i, feed := range feeds {
		fmt.Printf("%d. %s\n", i+1, feed.Name)
		fmt.Printf("   URL: %s\n", feed.Url)
		fmt.Printf("   Added by: %s\n", feed.UserName)
		fmt.Printf("   Added on: %s\n", feed.CreatedAt.Time.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}
	
	return nil
}

// handlerFollowFeed handles the follow command to follow an existing feed
func handlerFollowFeed(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("feed URL is required")
	}

	// Get the current user
	if s.Config == nil || s.Config.CurrentUserName == "" {
		return errors.New("you must be logged in to follow a feed")
	}

	ctx := context.Background()
	
	// Get the current user from the database
	user, err := s.Db.GetUserByName(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	
	feedURL := cmd.Args[0]
	
	// Get the feed by URL
	feed, err := s.Db.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("feed with URL %s not found: %w", feedURL, err)
	}
	
	// Create a new feed follow
	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	}
	
	feedFollow, err := s.Db.CreateFeedFollow(ctx, createFeedFollowParams)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"feed_follows_user_id_feed_id_key\" (SQLSTATE 23505)" {
			return fmt.Errorf("you are already following this feed")
		}
		return fmt.Errorf("failed to follow feed: %w", err)
	}
	
	fmt.Printf("You are now following feed \"%s\" as user %s\n", feedFollow.FeedName, feedFollow.UserName)
	
	return nil
}

// handlerListFollowing handles the following command to list feeds the user follows
func handlerListFollowing(s *cli.State, cmd cli.Command) error {
	// Get the current user
	if s.Config == nil || s.Config.CurrentUserName == "" {
		return errors.New("you must be logged in to see your followed feeds")
	}

	ctx := context.Background()
	
	// Get the current user from the database
	user, err := s.Db.GetUserByName(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	
	// Get the feeds user is following
	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get followed feeds: %w", err)
	}
	
	if len(feedFollows) == 0 {
		fmt.Printf("User %s is not following any feeds\n", user.Name)
		return nil
	}
	
	fmt.Printf("User %s is following %d feeds:\n\n", user.Name, len(feedFollows))
	
	// Display each followed feed
	for i, follow := range feedFollows {
		fmt.Printf("%d. %s\n", i+1, follow.FeedName)
		fmt.Printf("   URL: %s\n", follow.FeedUrl)
		fmt.Printf("   Following since: %s\n", follow.CreatedAt.Time.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}
	
	return nil
}