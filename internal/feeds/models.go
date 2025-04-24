package feeds

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/abahnj/rssagg/internal/database"
	"github.com/google/uuid"
)

// RSSFeed represents an RSS feed with its channels and items
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem represents a single item in an RSS feed
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Service handles feed operations
type Service struct {
	DB database.Queries
}

// NewService creates a new feed service
func NewService(db database.Queries) *Service {
	return &Service{
		DB: db,
	}
}

// FetchFeed retrieves and parses an RSS feed from the given URL
func (s *Service) FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
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

// CreateFeed adds a new feed to the database
func (s *Service) CreateFeed(ctx context.Context, name, url string, userID uuid.UUID) (database.Feed, error) {
	// Check if feed already exists
	existingFeed, err := s.DB.GetFeedByURL(ctx, url)
	if err == nil {
		// Feed already exists
		return existingFeed, nil
	}

	// Create a new feed record
	createFeedParams := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   name,
		Url:    url,
		UserID: userID,
	}

	feed, err := s.DB.CreateFeed(ctx, createFeedParams)
	if err != nil {
		return database.Feed{}, fmt.Errorf("failed to create feed: %w", err)
	}

	return feed, nil
}

// GetAllFeeds returns all feeds with their creator information
func (s *Service) GetAllFeeds(ctx context.Context) ([]database.GetFeedsWithUsersRow, error) {
	feeds, err := s.DB.GetFeedsWithUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get feeds: %w", err)
	}
	return feeds, nil
}

// FollowFeed creates a new feed follow
func (s *Service) FollowFeed(ctx context.Context, feedURL string, userID uuid.UUID) (database.CreateFeedFollowRow, error) {
	// Get the feed by URL
	feed, err := s.DB.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("feed with URL %s not found: %w", feedURL, err)
	}

	// Create a new feed follow
	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: userID,
		FeedID: feed.ID,
	}

	feedFollow, err := s.DB.CreateFeedFollow(ctx, createFeedFollowParams)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"feed_follows_user_id_feed_id_key\" (SQLSTATE 23505)" {
			return database.CreateFeedFollowRow{}, errors.New("you are already following this feed")
		}
		return database.CreateFeedFollowRow{}, fmt.Errorf("failed to follow feed: %w", err)
	}

	return feedFollow, nil
}

// GetFollowedFeeds returns all feeds a user is following
func (s *Service) GetFollowedFeeds(ctx context.Context, userID uuid.UUID) ([]database.GetFeedFollowsForUserRow, error) {
	follows, err := s.DB.GetFeedFollowsForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followed feeds: %w", err)
	}
	return follows, nil
}

// UnfollowFeed removes a feed follow for a user
func (s *Service) UnfollowFeed(ctx context.Context, feedURL string, userID uuid.UUID) error {
	params := database.DeleteFeedFollowParams{
		UserID: userID,
		Url:    feedURL,
	}
	
	err := s.DB.DeleteFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}
	
	return nil
}

// GetNextFeedToFetch gets the next feed that should be fetched
func (s *Service) GetNextFeedToFetch(ctx context.Context) (database.Feed, error) {
	feed, err := s.DB.GetNextFeedToFetch(ctx)
	if err != nil {
		return database.Feed{}, fmt.Errorf("failed to get next feed to fetch: %w", err)
	}
	return feed, nil
}

// MarkFeedFetched updates the last_fetched_at timestamp for a feed
func (s *Service) MarkFeedFetched(ctx context.Context, feedID uuid.UUID) error {
	err := s.DB.MarkFeedFetched(ctx, feedID)
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}
	return nil
}

// ScrapeFeeds fetches and processes a single feed
func (s *Service) ScrapeFeed(ctx context.Context) error {
	// Get the next feed to fetch
	feed, err := s.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}
	
	// Log which feed we're about to fetch
	fmt.Printf("\nFetching feed: %s (%s)\n", feed.Name, feed.Url)
	
	// Fetch the feed content
	rssFeed, err := s.FetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed content: %w", err)
	}
	
	// Process the feed items
	fmt.Printf("Found %d posts in feed\n\n", len(rssFeed.Channel.Item))
	
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("- %s\n  %s\n  Published: %s\n\n", 
			item.Title, 
			item.Link, 
			item.PubDate)
	}
	
	// Mark the feed as fetched
	err = s.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}
	
	return nil
}