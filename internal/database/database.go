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

	itemsChan := make(chan rss.Item, 100)
	errChan := make(chan error, 1)

	const numWorkers = 5
	for i := 0; i < numWorkers; i++ {
		go fetchRows(rows, itemsChan, errChan)
	}

	var items []rss.Item
	doneWorkers := 0
	for doneWorkers < numWorkers {
		select {
		case item := <-itemsChan:
			items = append(items, item)
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
			doneWorkers++
		}
	}

	return items, nil
}

// fetchRows reads rows from the *sql.Rows and sends rss.Item objects to itemsChan.
// It signals completion by sending nil to errChan.
func fetchRows(rows *sql.Rows, itemsChan chan<- rss.Item, errChan chan<- error) {
	for rows.Next() {
		var item rss.Item
		if err := rows.Scan(&item.Title, &item.Link, &item.Description); err != nil {
			errChan <- err
			return
		}
		itemsChan <- item
	}
	if err := rows.Err(); err != nil {
		errChan <- err
		return
	}
	errChan <- nil
}
