package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/talhag3/todoapp/internal/config"
	"github.com/talhag3/todoapp/internal/database"
	"github.com/talhag3/todoapp/internal/handler"
	"github.com/talhag3/todoapp/internal/pkg/pkg/logger"
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

	log := logger.New(logger.Config{
		Level:     "debug",
		Format:    "json", // "json" in prod, "text" while developing
		AddSource: true,
	})

	slog.SetDefault(log)

	// Connect to database
	pool, sqlDB, err := database.New(ctx, config.ConnectionString(conf))
	if err != nil {
		log.Error("failed to connect to database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	slog.Info("Connect to the Database")
	defer pool.Close()
	defer sqlDB.Close()

	if err := database.RunMigration(sqlDB, "migrations"); err != nil {
		slog.Error("Migration Failed ", slog.String("err", err.Error()))
		os.Exit(1)
	}

	taskRepo := repository.NewTaskRepository(pool, log)
	taskService := service.NewTaskService(taskRepo, log)

	taskH := handler.NewTaskHandler(taskService)

	// Fiber

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return response.Error(c, err)
		},
	})

	router.Register(router.Deps{
		App:   app,
		TaskH: taskH,
	})

	// Start Fiber server in a goroutine
	go func() {
		if err := app.Listen(":" + strconv.Itoa(conf.AppPort)); err != nil {
			log.Error("Failed to start Fiber server", slog.String("err", err.Error()))
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)

	// LISTEN FOR OS SIGNALS
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // Wait for signal
	cancel()  // Cancel context
	log.Info("Shutting down gracefully...")

	// Wait for Fiber to shutdown
	if err := app.Shutdown(); err != nil {
		log.Error("Failed to shutdown Fiber server:", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
