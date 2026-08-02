package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/talhag3/todoapp/internal/config"
	"github.com/talhag3/todoapp/internal/database"
	"github.com/talhag3/todoapp/internal/handler"
	"github.com/talhag3/todoapp/internal/pkg/response"
	"github.com/talhag3/todoapp/internal/repository"
	"github.com/talhag3/todoapp/internal/router"
	"github.com/talhag3/todoapp/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.LoadConf()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to database
	pool, sqlDB, err := database.New(ctx, config.ConnectionString(conf))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Connect to the Database")
	defer pool.Close()
	defer sqlDB.Close()

	if err := database.RunMigration(sqlDB, "migrations"); err != nil {
		log.Fatalf("Migration Failed %v", err)
		os.Exit(1)
	}

	taskRepo := repository.NewTaskRepository(pool)
	taskService := service.NewTaskService(taskRepo)

	taskH := handler.NewTaskHandler(taskService)

	// Fiber

	app := fiber.New(fiber.Config{
		ErrorHandler: response.ErrorHandler,
	})

	router.Register(router.Deps{
		App:   app,
		TaskH: taskH,
	})

	// Start Fiber server in a goroutine
	go func() {
		if err := app.Listen(":" + strconv.Itoa(conf.AppPort)); err != nil {
			log.Fatalf("Failed to start Fiber server: %v", err)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)

	// LISTEN FOR OS SIGNALS
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // Wait for signal
	cancel()  // Cancel context
	log.Println("Shutting down gracefully...")

	// Wait for Fiber to shutdown
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown Fiber server: %v", err)
	}
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
