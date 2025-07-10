package rss

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
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

func GetRSSFeeds(url string) (*RSS, error) {

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("status code isn't OK")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, err
	}

	return &rss, nil
}
