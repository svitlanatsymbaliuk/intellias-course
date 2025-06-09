package main

import (
	"fmt"

	"example.com/basics-practice/internal"
)

func main() {

	rss, err := internal.GetRSSFeeds("https://dou.ua/feed")

	if err != nil {
		fmt.Println("Error fetching RSS feeds:", err)
		return
	}

	fmt.Println("Feed Title:", rss.Channel.Title)
	for _, item := range rss.Channel.Items {
		fmt.Println("\n------------------------------------------------------------------------------------------")
		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Description:", item.Description)
	}

}
