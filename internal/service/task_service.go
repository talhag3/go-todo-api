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
	UpdateTask(ctx context.Context, id int64, title, description string, done bool) (*domain.Task, error)
	DeleteTask(ctx context.Context, id int64) error
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

func (svc *taskService) UpdateTask(ctx context.Context, id int64, title, description string, done bool) (*domain.Task, error) {

	task, err := svc.taskRepo.GetById(ctx, id)

	if err != nil {
		return nil, err
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

	err = svc.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (svc *taskService) DeleteTask(ctx context.Context, id int64) error {
	return svc.taskRepo.Delete(ctx, id)
}
