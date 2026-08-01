package handler

import "github.com/gofiber/fiber/v3"

type TaskHandler struct {
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) Check(c fiber.Ctx) error {

	type data struct {
		Message string
	}

	return c.JSON(data{Message: "Testing"})
}
