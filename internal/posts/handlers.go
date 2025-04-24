package posts

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abahnj/rssagg/internal/cli"
	"github.com/abahnj/rssagg/internal/database"
)

// HandlerBrowse handles the browse command to view posts from followed feeds
func HandlerBrowse(s *cli.State, cmd cli.Command, user database.User) error {
	ctx := context.Background()
	service := NewService(*s.Db)
	
	// Default limit is 10 posts
	limit := int32(10)
	
	// Parse limit argument if provided
	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit value: %w", err)
		}
		limit = int32(parsedLimit)
	}
	
	// Get posts for the user
	posts, err := service.GetPostsForUser(ctx, user.ID, limit)
	if err != nil {
		return err
	}
	
	if len(posts) == 0 {
		fmt.Println("No posts found in your followed feeds")
		return nil
	}
	
	fmt.Printf("Found %d posts from your followed feeds:\n\n", len(posts))
	
	// Display each post
	for i, post := range posts {
		fmt.Printf("%d. %s\n", i+1, post.Title)
		fmt.Printf("   Feed: %s\n", post.FeedName)
		fmt.Printf("   URL: %s\n", post.Url)
		
		// Show published date if available
		if post.PublishedAt.Valid {
			fmt.Printf("   Published: %s\n", post.PublishedAt.Time.Format("2006-01-02 15:04:05"))
		}
		
		// Show description if available
		if post.Description.Valid {
			// Truncate description if it's too long
			description := post.Description.String
			if len(description) > 100 {
				description = description[:100] + "..."
			}
			fmt.Printf("   Description: %s\n", description)
		}
		
		fmt.Println()
	}
	
	return nil
}