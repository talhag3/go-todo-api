package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/talhag3/todoapp/internal/domain"
)

const INSERT_TASK_QUERY = `
		INSERT INTO tasks (title, description, done_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
}

type taskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) TaskRepository {
	return &taskRepo{pool: pool}
}

func (r *taskRepo) Create(ctx context.Context, task *domain.Task) error {
	err := r.pool.QueryRow(ctx, INSERT_TASK_QUERY, task.Title, task.Description, task.DoneAt, task.CreatedAt, task.UpdatedAt).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}
