package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/rss"
)

func Connect(connectionData string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connectionData)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}

func InitializeDatabase(db *sql.DB) error {

	// Create table if not exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS rss_items (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		link TEXT NOT NULL,
		description TEXT NOT NULL
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return nil
}

func InsertRSSItem(db *sql.DB, items []rss.Item) error {

	insertQuery := `INSERT INTO rss_items (title, link, description) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`

	for _, item := range items {
		_, err := db.Exec(insertQuery, item.Title, item.Link, item.Description)
		if err != nil {
			return fmt.Errorf("error inserting item '%s': %v", item.Title, err)
		}

		fmt.Printf("Inserted item: %s\n", item.Title)
	}

	return nil
}

func GetAllRSSItems(db *sql.DB) ([]rss.Item, error) {

	rows, err := db.Query("SELECT title, link, description FROM rss_items")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []rss.Item
	for rows.Next() {
		var item rss.Item
		if err := rows.Scan(&item.Title, &item.Link, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
