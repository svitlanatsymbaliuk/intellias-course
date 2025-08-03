package rss

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Feed struct {
	response http.Response
}

func NewFeed(urlStr string) *Feed {
	parsedURL, _ := url.Parse(urlStr)
	return &Feed{
		response: http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("<rss><channel><title>Example RSS</title></channel></rss>")),
			Header:     http.Header{"Content-Type": []string{"application/rss+xml"}},
			Status:     "200 OK",
			Request:    &http.Request{Method: "GET", URL: parsedURL},
		},
	}
}

func (rssFeed *Feed) Get() (*RSS, error) {

	body, err := io.ReadAll(rssFeed.response.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, err
	}

	return &rss, nil
}
