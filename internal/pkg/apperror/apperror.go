package apperror

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	CodeNotFound     ErrorCode = "NOT_FOUND"
	CodeValidation   ErrorCode = "VALIDATION_ERROR"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeForbidden    ErrorCode = "FORBIDDEN"
	CodeConflict     ErrorCode = "CONFLICT"
	CodeBadRequest   ErrorCode = "BAD_REQUEST"
	CodeInternal     ErrorCode = "INTERNAL_ERROR"
)

var statusForCode = map[ErrorCode]int{
	CodeNotFound:     404,
	CodeValidation:   422,
	CodeUnauthorized: 401,
	CodeForbidden:    403,
	CodeConflict:     409,
	CodeBadRequest:   400,
	CodeInternal:     500,
}

type AppError struct {
	Code       ErrorCode
	Message    string
	StatusCode int
	Details    map[string]string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:       CodeNotFound,
		Message:    message,
		StatusCode: statusForCode[CodeNotFound],
	}
}

func Validation(message string, details map[string]string) *AppError {
	return &AppError{
		Code:       CodeValidation,
		Message:    message,
		StatusCode: statusForCode[CodeValidation],
		Details:    details,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       CodeUnauthorized,
		Message:    message,
		StatusCode: statusForCode[CodeUnauthorized],
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:       CodeForbidden,
		Message:    message,
		StatusCode: statusForCode[CodeForbidden],
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:       CodeConflict,
		Message:    message,
		StatusCode: statusForCode[CodeConflict],
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:       CodeBadRequest,
		Message:    message,
		StatusCode: statusForCode[CodeBadRequest],
	}
}

func Internal(err error) *AppError {
	return &AppError{
		Code:       CodeInternal,
		Message:    "something went wrong, please try again later",
		StatusCode: statusForCode[CodeInternal],
		Err:        err,
	}
}

// As is a small convenience wrapper around the standard library's
// errors.As. Given any error, it tells you: "is this (or does it wrap)
// an *AppError?" and hands you the typed value if so.

func As(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
