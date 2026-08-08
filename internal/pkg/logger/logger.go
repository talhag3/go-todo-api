package logger

import (
	"context"
	"log/slog"
	"os"
)

type Config struct {
	Level     string // "debug" | "info" | "warn" | "error"
	Format    string // "json" | "text"
	AddSource bool
}

func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)

	handlerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type ctxKey struct{}

func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return log
	}
	return slog.Default()
}

// WithRequestID is a convenience used by your Fiber middleware. It takes
// the base logger and returns a NEW logger that has request_id baked
// into every future call — you don't repeat "request_id", reqID on
// every single log line in the handler.
func WithRequestID(log *slog.Logger, requestID string) *slog.Logger {
	return log.With(slog.String("request_id", requestID))
}
