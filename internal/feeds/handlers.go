package feeds

import (
	"context"
	"errors"
	"fmt"

	"github.com/abahnj/rssagg/internal/cli"
	"github.com/abahnj/rssagg/internal/database"
)

// HandlerAggregator handles the agg command to fetch and display a feed
func HandlerAggregator(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	
	// URL of the feed to fetch
	feedURL := "https://www.wagslane.dev/index.xml"
	
	service := NewService(*s.Db)
	
	// Fetch the feed
	feed, err := service.FetchFeed(ctx, feedURL)
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

// HandlerAddFeed handles the addfeed command to add a new RSS feed
func HandlerAddFeed(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("both feed name and URL are required")
	}

	ctx := context.Background()
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	
	service := NewService(*s.Db)
	
	// Create or retrieve the feed
	feed, err := service.CreateFeed(ctx, feedName, feedURL, user.ID)
	if err != nil {
		return err
	}
	
	if feed.Name != feedName {
		fmt.Printf("Feed already exists, using existing feed: %s\n", feed.Name)
	} else {
		// Print out the new feed details
		fmt.Printf("Feed added successfully:\n")
		fmt.Printf("  ID: %s\n", feed.ID)
		fmt.Printf("  Name: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  Created: %s\n", feed.CreatedAt.Time.Format("2006-01-02 15:04:05"))
	}
	
	// Now create a feed follow for the current user
	followParams := database.CreateFeedFollowParams{
		ID:     feed.ID,
		UserID: user.ID,
		FeedID: feed.ID,
	}
	
	feedFollow, err := s.Db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}
	
	fmt.Printf("\nYou are now following this feed as user %s\n", feedFollow.UserName)
	
	return nil
}

// HandlerListFeeds handles the feeds command to list all feeds
func HandlerListFeeds(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	
	service := NewService(*s.Db)
	
	// Fetch all feeds with user information
	feeds, err := service.GetAllFeeds(ctx)
	if err != nil {
		return err
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

// HandlerFollowFeed handles the follow command to follow an existing feed
func HandlerFollowFeed(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("feed URL is required")
	}

	ctx := context.Background()
	feedURL := cmd.Args[0]
	
	service := NewService(*s.Db)
	
	// Follow the feed
	feedFollow, err := service.FollowFeed(ctx, feedURL, user.ID)
	if err != nil {
		return err
	}
	
	fmt.Printf("You are now following feed \"%s\" as user %s\n", feedFollow.FeedName, feedFollow.UserName)
	
	return nil
}

// HandlerUnfollowFeed handles the unfollow command to unfollow a feed
func HandlerUnfollowFeed(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("feed URL is required")
	}

	ctx := context.Background()
	feedURL := cmd.Args[0]
	
	service := NewService(*s.Db)
	
	// First check if the feed exists
	_, err := service.DB.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("feed with URL %s not found: %w", feedURL, err)
	}
	
	// Unfollow the feed
	err = service.UnfollowFeed(ctx, feedURL, user.ID)
	if err != nil {
		return err
	}
	
	fmt.Printf("You have unfollowed feed with URL: %s\n", feedURL)
	
	return nil
}

// HandlerListFollowing handles the following command to list feeds the user follows
func HandlerListFollowing(s *cli.State, cmd cli.Command, user database.User) error {
	ctx := context.Background()
	
	service := NewService(*s.Db)
	
	// Get the feeds user is following
	feedFollows, err := service.GetFollowedFeeds(ctx, user.ID)
	if err != nil {
		return err
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