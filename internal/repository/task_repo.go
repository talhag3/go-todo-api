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

const GET_TASK_BY_QUERY = `
		SELECT id, title, description, done_at, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetById(ctx context.Context, id int64) (*domain.Task, error)
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

func (r *taskRepo) GetById(ctx context.Context, id int64) (*domain.Task, error) {
	task := &domain.Task{}
	err := r.pool.QueryRow(ctx, GET_TASK_BY_QUERY, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DoneAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get task by ID: %w", err)
	}
	return task, nil
}
