package database

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/rss"
)

//go:embed *.sql
var embedMigrations embed.FS

type Database struct {
	connect *sql.DB
}

func NewDatabase(connectData string) *Database {
	db, err := sql.Open("postgres", connectData)
	if err != nil {
		return nil
	}
	return &Database{connect: db}
}

func (db *Database) Close() {
	if db.connect != nil {
		db.connect.Close()
	}
}

func (db *Database) GetConnection() *sql.DB {
	return db.connect
}

func (db *Database) Migration() error {
	goose.SetBaseFS(embedMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Use "." because goose will look for migrations in the embedded FS root
	err = goose.Up(db.connect, ".")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (db *Database) Initialize() error {

	// Create table if not exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS rss_items (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		link TEXT NOT NULL,
		description TEXT NOT NULL
	);`

	_, err := db.connect.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return nil
}

func (db *Database) InsertRSSItem(items []rss.Item) error {

	insertQuery := `INSERT INTO rss_items (title, link, description) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`

	for _, item := range items {
		_, err := db.connect.Exec(insertQuery, item.Title, item.Link, item.Description)
		if err != nil {
			return fmt.Errorf("error inserting item '%s': %v", item.Title, err)
		}

		fmt.Printf("Inserted item: %s\n", item.Title)
	}

	return nil
}

func (db *Database) GetAllRSSItems() ([]rss.Item, error) {

	rows, err := db.connect.Query("SELECT title, link, description FROM rss_items")
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
