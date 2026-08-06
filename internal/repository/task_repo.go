package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/talhag3/todoapp/internal/domain"
)

const insertTaskQuery = `
		INSERT INTO tasks (title, description, done_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

const getTaskByIDQuery = `
		SELECT id, title, description, done_at, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

const updateTaskQuery = `
		UPDATE tasks
		SET title = $1, description = $2, done_at = $3, updated_at = $4
		WHERE id = $5
	`

const GET_ALL_TASKS_QUERY = `
	SELECT id, title, description, done_at, created_at, updated_at
	FROM tasks
	ORDER BY created_at DESC
`

const deleteTaskQuery = `DELETE FROM tasks WHERE id = $1`

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetById(ctx context.Context, id int64) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]*domain.Task, error)
}

type taskRepo struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewTaskRepository(pool *pgxpool.Pool, log *slog.Logger) TaskRepository {
	return &taskRepo{pool: pool, log: log}
}

func (r *taskRepo) Create(ctx context.Context, task *domain.Task) error {
	start := time.Now()
	err := r.pool.QueryRow(ctx, insertTaskQuery,
		task.Title, task.Description, task.DoneAt, task.CreatedAt, task.UpdatedAt,
	).Scan(&task.ID)

	r.log.Debug("executed query",
		slog.String("query", "insert_task"),
		slog.Duration("duration", time.Since(start)),
		slog.Bool("success", err == nil),
	)

	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (r *taskRepo) GetById(ctx context.Context, id int64) (*domain.Task, error) {
	start := time.Now()
	task := &domain.Task{}
	err := r.pool.QueryRow(ctx, getTaskByIDQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DoneAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	r.log.Debug("executed query",
		slog.String("query", "get_task_by_id"),
		slog.Int64("task_id", id),
		slog.Duration("duration", time.Since(start)),
		slog.Bool("success", err == nil),
	)

	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return task, nil
}

func (r *taskRepo) Update(ctx context.Context, task *domain.Task) error {
	start := time.Now()

	tag, err := r.pool.Exec(
		ctx,
		updateTaskQuery,
		task.Title,
		task.Description,
		task.DoneAt,
		time.Now(),
		task.ID,
	)

	r.log.Debug("executed query",
		slog.String("query", "update_task"),
		slog.Int64("task_id", task.ID),
		slog.Duration("duration", time.Since(start)),
		slog.Bool("success", err == nil),
	)

	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update task: %w", ErrNoRowsAffected)
	}
	return nil
}

func (r *taskRepo) Delete(ctx context.Context, id int64) error {
	start := time.Now()

	tag, err := r.pool.Exec(ctx, deleteTaskQuery, id)

	r.log.Debug("executed query",
		slog.String("query", "delete_task"),
		slog.Int64("task_id", id),
		slog.Duration("duration", time.Since(start)),
		slog.Bool("success", err == nil),
	)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete task: %w", ErrNoRowsAffected)
	}
	return nil
}

func (r *taskRepo) GetAll(ctx context.Context) ([]*domain.Task, error) {
	start := time.Now()
	rows, err := r.pool.Query(ctx, GET_ALL_TASKS_QUERY)

	if err != nil {
		r.log.Debug("executed query",
			slog.String("query", "get_all_tasks"),
			slog.Duration("duration", time.Since(start)),
			slog.Bool("success", false),
		)
		return nil, fmt.Errorf("get all tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DoneAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	r.log.Debug("executed query",
		slog.String("query", "get_all_tasks"),
		slog.Int("row_count", len(tasks)),
		slog.Duration("duration", time.Since(start)),
		slog.Bool("success", true),
	)
	return tasks, nil
}
