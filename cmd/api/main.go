package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/talhag3/todoapp/internal/config"
	"github.com/talhag3/todoapp/internal/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.LoadConf()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.New(ctx, config.ConnectionString(conf))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)

	// LISTEN FOR OS SIGNALS
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// START A BACKGROUND WORKER (Goroutine)
	go func() {
		<-sigChan // Wait for signal
		cancel()  // Cancel context
		log.Println("Shutting down gracefully...")
	}()

	log.Println("Database connected successfully")
	<-ctx.Done() // Wait for context cancellation
}

/*
Shutting down gracefully

1. OS sends SIGINT
   ↓
2. signal.Notify puts SIGINT in sigChan
   ↓
3. Goroutine receives SIGINT from sigChan
   ↓
4. Goroutine calls cancel()
   ↓
5. ctx.Done() channel gets a value
   ↓
6. Main thread unblocks from <-ctx.Done()
   ↓
7. defer db.Close() runs (closes DB connections)
   ↓
8. App exits cleanly

*/
