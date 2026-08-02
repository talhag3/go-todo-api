package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is the canonical error type used across the app.
// Each instance carries: kind, user-facing message, HTTP status, and an optional wrapped cause.
type AppError struct {
	Kind    string // "not_found", "validation", "unauthorized", "conflict", "internal"
	Message string // safe to expose to client
	Status  int
	Cause   error // underlying error, never serialized
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

func New(kind, message string, status int, cause error) *AppError {
	return &AppError{Kind: kind, Message: message, Status: status, Cause: cause}
}

// Common constructors
func NotFound(msg string) *AppError   { return New("not_found", msg, http.StatusNotFound, nil) }
func Validation(msg string) *AppError { return New("validation", msg, http.StatusBadRequest, nil) }
func Unauthorized(msg string) *AppError {
	return New("unauthorized", msg, http.StatusUnauthorized, nil)
}
func Forbidden(msg string) *AppError { return New("forbidden", msg, http.StatusForbidden, nil) }
func Conflict(msg string) *AppError  { return New("conflict", msg, http.StatusConflict, nil) }
func Internal(msg string, cause error) *AppError {
	return New("internal", msg, http.StatusInternalServerError, cause)
}

// Map converts any error into an *AppError.
// Default: internal server error (avoid leaking details).
func Map(err error) *AppError {
	if err == nil {
		return nil
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	// Fallback: never expose raw error to clients.
	return Internal("something went wrong", err)
}
