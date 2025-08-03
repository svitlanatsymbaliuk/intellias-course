package database

import (
    "os"
    "testing"

    "github.com/svitlanatsymbaliuk/intellias-course/internal/rss"
)

func getTestDB() *Database {
    connStr := os.Getenv("TEST_DB_URL")
    if connStr == "" {
        connStr = "postgres://postgres:postgres@localhost:5432/go_app_db?sslmode=disable"
    }
    db := NewDatabase(connStr)
    return db
}

func TestNewDatabaseAndClose(t *testing.T) {
    db := getTestDB()
    if db == nil {
        t.Fatal("NewDatabase returned nil")
    }
    db.Close()
}

func TestMigration(t *testing.T) {
    db := getTestDB()
    if db == nil {
        t.Fatal("NewDatabase returned nil")
    }
    defer db.Close()
    if err := db.Migration(); err != nil {
        t.Fatalf("Migration failed: %v", err)
    }
}

func TestInitialize(t *testing.T) {
    db := getTestDB()
    if db == nil {
        t.Fatal("NewDatabase returned nil")
    }
    defer db.Close()
    if err := db.Initialize(); err != nil {
        t.Fatalf("Initialize failed: %v", err)
    }
}

func TestInsertAndGetAllRSSItems(t *testing.T) {
    db := getTestDB()
    if db == nil {
        t.Fatal("NewDatabase returned nil")
    }
    defer db.Close()
    if err := db.Initialize(); err != nil {
        t.Fatalf("Initialize failed: %v", err)
    }

    items := []rss.Item{
        {Title: "TestTitle1", Link: "http://test1", Description: "Desc1"},
        {Title: "TestTitle2", Link: "http://test2", Description: "Desc2"},
    }
    if err := db.InsertRSSItem(items); err != nil {
        t.Fatalf("InsertRSSItem failed: %v", err)
    }

    got, err := db.GetAllRSSItems()
    if err != nil {
        t.Fatalf("GetAllRSSItems failed: %v", err)
    }
    if len(got) < 2 {
        t.Errorf("Expected at least 2 items,