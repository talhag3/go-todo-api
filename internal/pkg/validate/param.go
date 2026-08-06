// pkg/validate/param.go
//
// Helpers for validating things that come from the URL itself — path
// params and query params — as opposed to validate.Struct(), which
// validates the JSON body. Different source, same job: turn a raw
// string into a typed value, or a proper *AppError if it can't.s
package validate

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/talhag3/todoapp/internal/pkg/apperror"
)

// ParamInt64 reads c.Params(name), parses it as int64, and returns a
// ready-to-use *AppError on failure. `field` is the JSON key used in
// the error details map — so a caller validating "id" gets
// {"id": "..."} back, and a caller validating "task_id" (if you ever
// nest routes like /users/:user_id/tasks/:task_id) gets that instead.
//
// func signature convention: Go often names the "thing being parsed"
// in the function name (ParamInt64) rather than a vague generic verb —
// makes call sites self-documenting without needing to open the file.
func ParamInt64(c fiber.Ctx, field string) (int64, error) {
	raw := c.Params(field)
	if raw == "" {
		return 0, apperror.Validation(field+" is required", map[string]string{
			field: field + " is required",
		})
	}

	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, apperror.Validation("invalid "+field, map[string]string{
			field: field + " must be a valid integer",
		})
	}

	return val, nil
}
