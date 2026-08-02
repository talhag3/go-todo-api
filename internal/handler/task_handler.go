package handler

import (
	"fmt"
	"strconv"
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

type editTaskReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DoneAt      *bool  `json:"done_at"` // Accepts true/false
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

func (h *TaskHandler) GetTask(c fiber.Ctx) error {

	taskIDStr := c.Params("id")
	if taskIDStr == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Task ID is required",
			"errors":  map[string][]string{"id": {"Task ID is required"}},
		})
	}

	// Convert taskID from string to int64
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Invalid Task ID",
			"errors":  map[string][]string{"id": {"Task ID must be a valid integer"}},
		})
	}

	task, err := h.svc.GetTaskByID(c, taskID)

	if err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"success": false,
			"message": "Task not found",
			"errors":  map[string][]string{"id": {"Task with the given ID does not exist"}},
		})
	}

	// Return the task if found
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task retrieved successfully",
		"data":    task,
	})
}

func (h *TaskHandler) GetTasks(c fiber.Ctx) error {

	tasks, err := h.svc.GetAllTasks(c)

	if err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Issue",
			"errors":  err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Tasks retrieved successfully",
		"data":    tasks,
	})
}

func (h *TaskHandler) EditTask(c fiber.Ctx) error {

	taskIDStr := c.Params("id")
	if taskIDStr == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Task ID is required",
			"errors":  map[string][]string{"id": {"Task ID is required"}},
		})
	}

	// Convert taskID from string to int64
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Invalid Task ID",
			"errors":  map[string][]string{"id": {"Task ID must be a valid integer"}},
		})
	}

	req := editTaskReq{}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"errors":  map[string][]string{"body": {err.Error()}},
		})
	}

	// Validate the request
	validationErrors := validateEditTaskReq(req)
	if len(validationErrors) > 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
	}

	// Handle DoneAt: Default to false if not provided
	doneAt := false
	if req.DoneAt != nil {
		fmt.Println(doneAt)
		doneAt = *req.DoneAt
	}

	updatedTask, err := h.svc.UpdateTask(c, taskID, req.Title, req.Description, doneAt)

	if err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"success": false,
			"message": "Task not found",
			"errors":  map[string][]string{"id": {"Task with the given ID does not exist"}},
		})
	}

	// Return the updated task
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task updated successfully",
		"data":    updatedTask,
	})
}

func (h *TaskHandler) DeleteTask(c fiber.Ctx) error {
	// Extract taskID from URL
	taskIDStr := c.Params("id")
	if taskIDStr == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Task ID is required",
			"errors":  map[string][]string{"id": {"Task ID is required"}},
		})
	}

	// Convert taskID to int64
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"message": "Invalid Task ID",
			"errors":  map[string][]string{"id": {"Task ID must be a valid integer"}},
		})
	}

	// Call the service to delete the task
	_ = h.svc.DeleteTask(c.Context(), taskID)

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task deleted successfully",
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

func validateEditTaskReq(req editTaskReq) map[string][]string {
	errors := make(map[string][]string)

	// Title is required
	if req.Title == "" {
		errors["title"] = append(errors["title"], "Title is required")
	} else if len(req.Title) < 3 || len(req.Title) > 100 {
		errors["title"] = append(errors["title"], "Title must be between 3 and 100 characters")
	}

	// Description is optional but must meet length requirements if provided
	if req.Description != "" && (len(req.Description) < 5 || len(req.Description) > 500) {
		errors["description"] = append(errors["description"], "Description must be between 5 and 500 characters")
	}

	return errors
}
