package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/talhag3/todoapp/internal/domain"
	"github.com/talhag3/todoapp/internal/pkg/apperror"
	"github.com/talhag3/todoapp/internal/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, title, description string) (*domain.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*domain.Task, error)
	UpdateTask(ctx context.Context, id int64, title, description string, done bool) (*domain.Task, error)
	DeleteTask(ctx context.Context, id int64) error
	GetAllTasks(ctx context.Context) ([]*domain.Task, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
	log      *slog.Logger
}

func NewTaskService(repo repository.TaskRepository, log *slog.Logger) TaskService {
	return &taskService{taskRepo: repo, log: log}
}

func (svc *taskService) CreateTask(ctx context.Context, title, description string) (*domain.Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, apperror.Validation("validation failed", map[string]string{
			"title": "this field is required",
		})
	}

	now := time.Now()
	task := &domain.Task{
		Title:       title,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := svc.taskRepo.Create(ctx, task)

	if err != nil {
		svc.log.Error("failed to create task",
			slog.String("title", title),
			slog.String("err", err.Error()),
		)
		return nil, apperror.Internal(err)
	}

	svc.log.Info("task created", slog.Int64("task_id", task.ID))
	return task, nil
}

func (svc *taskService) GetTaskByID(ctx context.Context, id int64) (*domain.Task, error) {
	// log.With() returns a CHILD logger carrying task_id on every call
	// made through it from here on — same idea as WithRequestID from
	// the middleware, just scoped to this one method instead of a
	// whole request. Saves repeating slog.Int64("task_id", id) on
	// every single log line below.
	log := svc.log.With(slog.Int64("task_id", id))

	task, err := svc.taskRepo.GetById(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("task not found")
			return nil, apperror.NotFound("task not found")
		}
		log.Error("failed to get task", slog.String("err", err.Error()))
		return nil, apperror.Internal(err)
	}

	return task, nil
}

func (svc *taskService) UpdateTask(ctx context.Context, id int64, title, description string, done bool) (*domain.Task, error) {
	log := svc.log.With(slog.Int64("task_id", id))

	title = strings.TrimSpace(title)
	if title == "" {
		return nil, apperror.Validation("validation failed", map[string]string{
			"title": "this field is required",
		})
	}

	task, err := svc.taskRepo.GetById(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("task not found for update")
			return nil, apperror.NotFound("task not found")
		}
		log.Error("failed to load task for update", slog.String("err", err.Error()))
		return nil, apperror.Internal(err)
	}

	task.Title = title
	task.Description = description
	task.UpdatedAt = time.Now()

	if done {
		now := time.Now()
		task.DoneAt = &now
	} else {
		task.DoneAt = nil
	}

	if err := svc.taskRepo.Update(ctx, task); err != nil {
		if errors.Is(err, repository.ErrNoRowsAffected) {
			log.Warn("task not found for update")
			return nil, apperror.NotFound("task not found")
		}
		log.Error("failed to update task", slog.String("err", err.Error()))
		return nil, apperror.Internal(err)
	}
	return task, nil
}

func (svc *taskService) DeleteTask(ctx context.Context, id int64) error {
	log := svc.log.With(slog.Int64("task_id", id))

	if err := svc.taskRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNoRowsAffected) {
			log.Warn("task not found for delete")
			return apperror.NotFound("task not found")
		}
		log.Error("failed to delete task", slog.String("err", err.Error()))
		return apperror.Internal(err)
	}

	log.Info("task deleted")
	return nil
}

func (svc *taskService) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	tasks, err := svc.taskRepo.GetAll(ctx)
	if err != nil {
		svc.log.Error("failed to list tasks", slog.String("err", err.Error()))
		return nil, apperror.Internal(err)
	}
	return tasks, nil
}
