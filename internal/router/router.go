package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/talhag3/todoapp/internal/handler"
)

type Deps struct {
	App   *fiber.App
	TaskH *handler.TaskHandler
}

func Register(d Deps) {
	app := d.App

	api := app.Group("/api/v1")

	tasks := api.Group("/tasks")

	tasks.Get("/test", d.TaskH.Check)
}
