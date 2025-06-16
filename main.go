package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/svitlanatsymbaliuk/intellias-course/internal"
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

	dbHost := "db"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := "go_app_db"

	connectionData := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connectionData)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
	} else {
		fmt.Println("Successfully connected to the database!")
	}

	// Create table if not exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS rss_items (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		link TEXT NOT NULL,
		description TEXT NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// Insert RSS items into the table
	insertQuery := `INSERT INTO rss_items (title, link, description) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`
	for _, item := range rss.Channel.Items {
		_, err := db.Exec(insertQuery, item.Title, item.Link, item.Description)
		if err != nil {
			fmt.Printf("Error inserting item '%s': %v\n", item.Title, err)
		} else {
			// Check if the item was inserted
			var count int
			checkQuery := `SELECT COUNT(*) FROM rss_items WHERE title = $1 AND link = $2`
			err = db.QueryRow(checkQuery, item.Title, item.Link).Scan(&count)
			if err != nil {
				fmt.Printf("Error checking item '%s': %v\n", item.Title, err)
			} else if count > 0 {
				fmt.Printf("Item '%s' successfully inserted.\n", item.Title)
			} else {
				fmt.Printf("Item '%s' was not inserted.\n", item.Title)
			}
		}
	}
}
