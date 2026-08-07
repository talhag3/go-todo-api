package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/talhag3/todoapp/internal/dto"
	"github.com/talhag3/todoapp/internal/pkg/apperror"
	"github.com/talhag3/todoapp/internal/pkg/response"
	"github.com/talhag3/todoapp/internal/pkg/validate"
	"github.com/talhag3/todoapp/internal/service"
)

type TaskHandler struct {
	svc service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Create(c fiber.Ctx) error {
	req := dto.TaskCreateRequest{}

	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(c, err) // *AppError already, response.Error handles it directly
	}

	newTask, err := h.svc.CreateTask(c, req.Title, req.Description)

	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, newTask)
}

func (h *TaskHandler) GetTask(c fiber.Ctx) error {

	taskID, err := validate.ParamInt64(c, "id")
	if err != nil {
		return response.Error(c, err)
	}

	task, err := h.svc.GetTaskByID(c.Context(), taskID)
	if err != nil {

		return response.Error(c, err)
	}

	return response.OK(c, task)
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

	taskID, err := validate.ParamInt64(c, "id")
	if err != nil {
		return response.Error(c, err)
	}

	req := dto.TaskEditRequest{}
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(c, err)
	}

	doneAt := false
	if req.DoneAt != nil {
		doneAt = *req.DoneAt
	}

	updatedTask, err := h.svc.UpdateTask(c, taskID, req.Title, req.Description, doneAt)

	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, updatedTask)
}

func (h *TaskHandler) DeleteTask(c fiber.Ctx) error {

	taskID, err := validate.ParamInt64(c, "id")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.DeleteTask(c.Context(), taskID); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)

}
