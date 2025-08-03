package rss

import (
	"testing"
)

func TestNewFeed(t *testing.T) {
	url := "http://example.com/rss"
	feed := NewFeed(url)
	if feed == nil {
		t.Fatal("NewFeed returned nil")
	}
	if feed.response.StatusCode != 200 {
		t.Errorf("Expected StatusCode 200, got %d", feed.response.StatusCode)
	}
	if feed.response.Header.Get("Content-Type") != "application/rss+xml" {
		t.Errorf("Expected Content-Type application/rss+xml, got %s", feed.response.Header.Get("Content-Type"))
	}
}

func TestFeed_Get(t *testing.T) {
	url := "http://example.com/rss"
	feed := NewFeed(url)
	rss, err := feed.Get()
	if err != nil {
		t.Fatalf("Feed.Get returned error: %v", err)
	}
	if rss == nil {
		t.Fatal("Feed.Get returned nil RSS")
	}
	if rss.Channel.Title != "Example RSS" {
		t.Errorf("Expected Channel.Title 'Example RSS', got '%s'", rss.Channel.Title)
	}
}
