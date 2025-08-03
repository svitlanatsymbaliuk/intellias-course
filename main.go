package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/config"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/database"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/rss"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/server"
)

func main() {

	fmt.Println("Starting RSS Feed Application...")
	url := "http://podcast.dou.ua/rss"

	rssFeed := rss.NewFeed(url)
	rss, err := rssFeed.Get()
	if err != nil {
		fmt.Printf("error fetching RSS feeds: %v", err)
		return
	}
	fmt.Println("RSS Feed fetched successfully.")

	fmt.Println("Connecting to the database...")
	connectionData := config.NewConnectionDB().URL
	database := database.NewDatabase(connectionData)
	if database == nil {
		fmt.Printf("error connecting to the database: %v", err)
		return
	}
	fmt.Println("Database connection established.")

	defer database.Close()

	if err := database.Migration(); err != nil {
		fmt.Printf("error applying migrations: %v", err)
		return
	}
	fmt.Println("Database migrations applied successfully.")

	if err := database.Initialize(); err != nil {
		fmt.Printf("error initializing database: %v", err)
		return
	}

	if err := database.InsertRSSItem(rss.Channel.Items); err != nil {
		fmt.Printf("error inserting RSS items: %v", err)
		return
	}
	fmt.Println("RSS items inserted into the database successfully.")

	// Start the REST API server - ToDo: move to a separate package
	e := server.New()
	e.GET("/rss", func(c echo.Context) error {
		items, err := database.GetAllRSSItems()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch items"})
		}
		return c.JSON(http.StatusOK, items)
	})

	fmt.Println("REST API running on :8080")
	if err := e.Start(":8080"); err != nil {
		e.Logger.Info("Shutting down the server")
	}

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server forced to shutdown: ", err)
	}
}
