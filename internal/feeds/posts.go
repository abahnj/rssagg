package feeds

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abahnj/rssagg/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreatePost adds a new post to the database
func (s *Service) CreatePost(ctx context.Context, feed database.Feed, item RSSItem) error {
	// Skip items with missing title or URL
	if item.Title == "" || item.Link == "" {
		return errors.New("post missing required title or URL")
	}

	// Parse the published date
	var publishedAt pgtype.Timestamp
	if item.PubDate != "" {
		parsedTime, err := parseRSSTime(item.PubDate)
		if err == nil {
			publishedAt.Time = parsedTime
			publishedAt.Valid = true
		}
	}

	// Prepare description
	var description pgtype.Text
	if item.Description != "" {
		description.String = item.Description
		description.Valid = true
	}

	// Create post parameters
	params := database.CreatePostParams{
		ID:          uuid.New(),
		Title:       item.Title,
		Url:         item.Link,
		Description: description,
		PublishedAt: publishedAt,
		FeedID:      feed.ID,
	}

	// Insert post into database
	_, err := s.DB.CreatePost(ctx, params)
	if err != nil {
		// If the error is a duplicate key error, just ignore it
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetPostsForUser fetches posts for a specific user with a limit
func (s *Service) GetPostsForUser(ctx context.Context, userID uuid.UUID, limit int32) ([]database.GetPostsForUserRow, error) {
	params := database.GetPostsForUserParams{
		UserID: userID,
		Limit:  limit,
	}

	posts, err := s.DB.GetPostsForUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts for user: %w", err)
	}

	return posts, nil
}

// parseRSSTime attempts to parse a time string from an RSS feed in various formats
func parseRSSTime(timeStr string) (time.Time, error) {
	// Common time formats found in RSS feeds
	formats := []string{
		time.RFC1123Z,     // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,      // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC3339,      // "2006-01-02T15:04:05Z07:00"
		time.RFC822Z,      // "02 Jan 06 15:04 -0700"
		time.RFC822,       // "02 Jan 06 15:04 MST"
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700",
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	// If we got here, none of the formats matched
	return time.Time{}, fmt.Errorf("could not parse time string: %s, last error: %w", timeStr, lastErr)
}