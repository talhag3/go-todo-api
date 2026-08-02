package response

import (
	"github.com/gofiber/fiber/v3"
	"github.com/talhag3/todoapp/internal/pkg/apperror"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type Meta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
	Pages   int `json:"pages"`
}

type ErrorBody struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func OK(c fiber.Ctx, data interface{}) error {
	return c.JSON(Response{Success: true, Data: data})
}

func OKWithMeta(c fiber.Ctx, data interface{}, meta *Meta) error {
	return c.JSON(Response{Success: true, Data: data, Meta: meta})
}

func Created(c fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{Success: true, Data: data})
}

func NoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func ErrorHandler(c fiber.Ctx, err error) error {
	ae := apperror.Map(err)

	return c.Status(ae.Status).JSON(Response{
		Success: false,
		Error: &ErrorBody{
			Type:    ae.Kind,
			Message: ae.Message,
		},
	})
}
