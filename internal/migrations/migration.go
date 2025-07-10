package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Run(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Use "." because goose will look for migrations in the embedded FS root
	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
