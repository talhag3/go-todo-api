package service

import (
	"context"
	"time"

	"github.com/talhag3/todoapp/internal/domain"
	"github.com/talhag3/todoapp/internal/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, title, description string) (*domain.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*domain.Task, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: repo}
}

func (svc *taskService) CreateTask(ctx context.Context, title, description string) (*domain.Task, error) {
	now := time.Now()
	task := &domain.Task{
		Title:       title,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := svc.taskRepo.Create(ctx, task)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (svc *taskService) GetTaskByID(ctx context.Context, id int64) (*domain.Task, error) {
	return svc.taskRepo.GetById(ctx, id)
}
