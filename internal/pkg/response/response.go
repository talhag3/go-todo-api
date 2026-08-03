package response

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/talhag3/todoapp/internal/pkg/apperror"
)

type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type successBody struct {
	Success bool `json:"success"`
	Data    any  `json:"data,ommitempty"`
	Meta    any  `json:"meta,ommitempty"`
}

type errorBody struct {
	Success bool      `json:"success"`
	Error   errorInfo `json:"error"`
}

type errorInfo struct {
	Code    apperror.ErrorCode `json:"code"`
	Message string             `json:"message"`
	Details map[string]string  `json:"details,omitempty"`
}

func OK(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(successBody{
		Success: true,
		Data:    data,
	})
}

func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(successBody{
		Success: true,
		Data:    data,
	})
}

func NoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func Paginated(c fiber.Ctx, data any, meta Meta) error {
	return c.Status(fiber.StatusOK).JSON(successBody{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

func Error(c fiber.Ctx, err error) error {
	// 1. Is it one of ours?
	if appErr, ok := apperror.As(err); ok {
		return sendError(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
	}

	// 2. Is it a *fiber.Error? This is what Fiber itself raises for
	// unmatched routes (404), wrong HTTP method (405), bad request body
	// parsing, payload too large, etc. Fiber's Error struct just carries
	// a Code (int) and Message (string) — no fields for our ErrorCode enum,
	// so we translate it ourselves.
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return sendError(c, fiberErr.Code, codeForStatus(fiberErr.Code), fiberErr.Message, nil)
	}

	// 3. Anything else is unexpected — treat as 500, don't leak details.
	appErr := apperror.Internal(err)
	return sendError(c, appErr.StatusCode, appErr.Code, appErr.Message, nil)
}

// sendError is a small private helper (lowercase = unexported, only
// visible inside this package) so the three branches above don't repeat
// the JSON-building code.
func sendError(c fiber.Ctx, status int, code apperror.ErrorCode, message string, details map[string]string) error {
	return c.Status(status).JSON(errorBody{
		Success: false,
		Error: errorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// codeForStatus gives Fiber's raw numeric statuses one of our named
// ErrorCodes, so the JSON response stays consistent even for errors
// Fiber generated internally rather than us.
func codeForStatus(status int) apperror.ErrorCode {
	switch status {
	case fiber.StatusNotFound:
		return apperror.CodeNotFound
	case fiber.StatusMethodNotAllowed:
		return apperror.CodeBadRequest
	case fiber.StatusUnauthorized:
		return apperror.CodeUnauthorized
	case fiber.StatusForbidden:
		return apperror.CodeForbidden
	default:
		return apperror.CodeInternal
	}
}
