package dto

import "time"

type TaskCreateRequest struct {
	Title       string     `json:"title" validate:"required,min=3,max=100"`
	Description string     `json:"description" validate:"omitempty,min=5,max=500"`
	DoneAt      *time.Time `json:"done_at"`
}

type TaskEditRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,min=5,max=500"`
	DoneAt      *bool  `json:"done_at"`
}
