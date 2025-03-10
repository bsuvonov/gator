package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
)


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



func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed

	xml.Unmarshal(body, &feed)

	p := bluemonday.UGCPolicy()
    feed.Channel.Description = p.Sanitize(html.UnescapeString(feed.Channel.Description))
	feed.Channel.Title = p.Sanitize(html.UnescapeString(feed.Channel.Title))

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = p.Sanitize(html.UnescapeString(feed.Channel.Item[i].Title))
		feed.Channel.Item[i].Description = p.Sanitize(html.UnescapeString(feed.Channel.Item[i].Description))
	}

	return &feed, nil
}