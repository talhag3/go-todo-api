package domain

import "time"

type Task struct {
	ID          int64
	Title       string
	Description string
	DoneAt      *time.Time // Need nil to represent SQL NULL
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
