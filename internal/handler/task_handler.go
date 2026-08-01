package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/talhag3/todoapp/internal/service"
)

type TaskHandler struct {
	svc service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

type taskReq struct {
	Title       string
	Description string
	DoneAt      *time.Time
}

func (h *TaskHandler) Create(c fiber.Ctx) error {
	req := taskReq{}

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"errors":  map[string][]string{"body": {err.Error()}},
		})
	}

	// Validate the request
	validationErrors := validateTaskReq(req)
	if len(validationErrors) > 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	// Create the Task on the DB

	newTask, err := h.svc.CreateTask(c, req.Title, req.Description)

	if err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"success": false,
			"message": "Internal Error",
			"errors":  map[string][]string{"body": {err.Error()}},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Task created successfully",
		"data":    newTask,
	})
}

func validateTaskReq(req taskReq) map[string][]string {
	errors := make(map[string][]string)

	// Validate Title
	if req.Title == "" {
		errors["title"] = append(errors["title"], "Title is required")
	} else if len(req.Title) < 3 || len(req.Title) > 100 {
		errors["title"] = append(errors["title"], "Title must be between 3 and 100 characters")
	}

	// Validate Description
	if req.Description != "" && (len(req.Description) < 5 || len(req.Description) > 500) {
		errors["description"] = append(errors["description"], "Description must be between 5 and 500 characters")
	}

	return errors
}
