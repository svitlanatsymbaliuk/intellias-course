package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/config"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/database"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/migrations"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/rss"
	"github.com/svitlanatsymbaliuk/intellias-course/internal/server"
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

	e := server.New()
	server.Get(e, db, "/rss")

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
