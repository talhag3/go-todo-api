package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/talhag3/todoapp/internal/pkg/logger"
)

func RequestLogger(base *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		requestId := uuid.NewString()
		reqLog := logger.WithRequestID(base, requestId)

		c.SetContext(logger.WithContext(c.Context(), reqLog))
		c.Set("X-Request-ID", requestId)

		err := c.Next() // runs the matched handler

		status := c.Response().StatusCode()

		level := slog.LevelInfo

		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		reqLog.Log(c.Context(), level, "request completed",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
		)

		return err
	}
}
