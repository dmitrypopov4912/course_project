package task_services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"task-manager/internal/dto"
	"task-manager/internal/storage/postgres_storage"
	"time"
)

var (
	ErrEmptyTitle      = errors.New("title cannot be empty")
	ErrInvalidPriority = errors.New("invalid priority value")
	ErrInvalidStatus   = errors.New("invalid status value")
	ErrDueDateInPast   = errors.New("due date cannot be in the past")
)

func validateTaskRequest(taskReq dto.TaskRequest) error {
	if strings.TrimSpace(taskReq.Title) == "" {
		return ErrEmptyTitle
	}
	if taskReq.Priority < 0 || taskReq.Priority > 2 {
		return ErrInvalidPriority
	}
	if taskReq.Status < 0 || taskReq.Status > 2 {
		return ErrInvalidStatus
	}
	if taskReq.DueDate != nil && taskReq.DueDate.Before(time.Now()) {
		return ErrDueDateInPast
	}
	return nil
}

type TaskService struct {
	storage *postgres_storage.Storage
}

func NewService(storage *postgres_storage.Storage) *TaskService {
	return &TaskService{storage: storage}
}

func (ts *TaskService) CreateTask(ctx context.Context, taskReq dto.TaskRequest) error {
	if err := validateTaskRequest(taskReq); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	task := convertTask(taskReq)
	if err := ts.storage.CreateTask(ctx, task); err != nil {
		return fmt.Errorf("ts.storage.CreateTask: %w", err)
	}
	return nil
}

func (ts *TaskService) GetTask(ctx context.Context, id int64) (*dto.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task id: %d", id)
	}
	task, err := ts.storage.GetTask(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ts.storage.GetTask: %w", err)
	}
	return task, nil
}

func (ts *TaskService) UpdateTask(ctx context.Context, taskReq dto.TaskRequest, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid task id: %d", id)
	}
	if err := validateTaskRequest(taskReq); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	task := convertTask(taskReq)
	task.ID = id
	if err := ts.storage.UpdateTask(ctx, task); err != nil {
		return fmt.Errorf("ts.storage.UpdateTask: %w", err)
	}
	return nil
}

func (ts *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid task id: %d", id)
	}
	if err := ts.storage.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("ts.storage.DeleteTask: %w", err)
	}
	return nil
}

func convertTask(taskReq dto.TaskRequest) dto.Task {
	return dto.Task{
		Title:       taskReq.Title,
		Description: taskReq.Description,
		DueDate:     taskReq.DueDate,
		Priority:    toPriority(taskReq.Priority),
		Status:      toStatus(taskReq.Status),
	}
}

func toPriority(priority int) dto.Priority {
	switch priority {
	case 0:
		return dto.PriorityLow
	case 1:
		return dto.ProrityMedium
	case 2:
		return dto.PriorityHigh
	default:
		return dto.PriorityNotSpecified
	}
}

func toStatus(status int) dto.Status {
	switch status {
	case 0:
		return dto.StatusTodo
	case 1:
		return dto.StatusWip
	case 2:
		return dto.StatusDone
	default:
		return dto.StatusNotSpecified
	}
}
