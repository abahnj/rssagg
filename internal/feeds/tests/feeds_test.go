package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abahnj/rssagg/internal/feeds"
)

func TestFetchFeed(t *testing.T) {
	t.Run("Successfully parse RSS feed", func(t *testing.T) {
		// Create a test server with a mock RSS feed
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check User-Agent header
			if r.Header.Get("User-Agent") != "gator" {
				t.Errorf("Expected User-Agent to be 'gator', got %s", r.Header.Get("User-Agent"))
			}
			
			// Return a simple RSS feed
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>
  <title>Test Feed</title>
  <link>https://example.com</link>
  <description>A test RSS feed &amp; more</description>
  <item>
    <title>Test Item "Quoted"</title>
    <link>https://example.com/item1</link>
    <description>Test description &lt;b&gt;with HTML&lt;/b&gt;</description>
    <pubDate>Mon, 01 Jan 2024 12:00:00 GMT</pubDate>
  </item>
</channel>
</rss>`))
		}))
		defer server.Close()
		
		// Create a service with a mock DB
		service := &feeds.Service{}
		
		// Fetch the feed
		ctx := context.Background()
		feed, err := service.FetchFeed(ctx, server.URL)
		
		// Check for errors
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Check feed values
		if feed.Channel.Title != "Test Feed" {
			t.Errorf("Expected title to be 'Test Feed', got %s", feed.Channel.Title)
		}
		
		if feed.Channel.Description != "A test RSS feed & more" {
			t.Errorf("Expected unescaped description, got %s", feed.Channel.Description)
		}
		
		// Check item values
		if len(feed.Channel.Item) != 1 {
			t.Fatalf("Expected 1 item, got %d", len(feed.Channel.Item))
		}
		
		item := feed.Channel.Item[0]
		if item.Title != `Test Item "Quoted"` {
			t.Errorf("Expected unescaped title, got %s", item.Title)
		}
		
		if item.Description != "Test description <b>with HTML</b>" {
			t.Errorf("Expected unescaped description, got %s", item.Description)
		}
	})
	
	t.Run("Handle error status code", func(t *testing.T) {
		// Create a test server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()
		
		// Create a service with a mock DB
		service := &feeds.Service{}
		
		// Fetch the feed
		ctx := context.Background()
		_, err := service.FetchFeed(ctx, server.URL)
		
		// Check for error
		if err == nil {
			t.Fatalf("Expected error for status 404, got nil")
		}
	})
}