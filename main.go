package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/config"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/database"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/migrations"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/rss"
)

func main() {

	url := "http://podcast.dou.ua/rss"

	rss, err := rss.GetRSSFeeds(url)
	if err != nil {
		fmt.Printf("error fetching RSS feeds: %v", err)
		return
	}

	connectionData := config.NewConnectionDB().URL
	db, err := database.Connect(connectionData)
	if err != nil {
		fmt.Printf("error connecting to the database: %v", err)
		return
	}

	if err := migrations.Run(db); err != nil {
		fmt.Printf("error applying migrations: %v", err)
		return
	}

	if err := database.InitializeDatabase(db); err != nil {
		fmt.Printf("error initializing database: %v", err)
		return
	}

	if err := database.InsertRSSItem(db, rss.Channel.Items); err != nil {
		fmt.Printf("error inserting RSS items: %v", err)
		return
	}

	defer db.Close()

	// Add REST API with Echo
	e := echo.New()
	e.GET("/rss", func(context echo.Context) error {
		items, err := database.GetAllRSSItems(db)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch items"})
		}
		return context.JSON(http.StatusOK, items)
	})

	fmt.Println("REST API running on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
