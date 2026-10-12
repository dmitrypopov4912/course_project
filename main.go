package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrEmptyTitle   = errors.New("title cannot be empty")
)

type Task struct {
	ID          int
	Title       string
	Description string
	Priority    int
	Completed   bool
	AddedAt     time.Time
	Deadline    time.Time
}

type Storage interface {
	CreateTask(task Task) (Task, error)
	GetTask(id int) (Task, error)
	UpdateTask(task Task) error
	DeleteTask(id int) error
}

type MemoryStorage struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
}

func NewMemoryStogare() *MemoryStorage {
	return &MemoryStorage{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (s *MemoryStorage) CreateTask(task Task) (Task, error) {
	if task.Title == "" {
		return Task{}, ErrEmptyTitle
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	task.ID = s.nextID
	s.nextID++
	task.AddedAt = time.Now()
	s.tasks[task.ID] = task
	return task, nil
}

func (s *MemoryStorage) GetTask(id int) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *MemoryStorage) UpdateTask(task Task) error {
	if task.Title == "" {
		return ErrEmptyTitle
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}
	s.tasks[task.ID] = task
	return nil
}

func (s *MemoryStorage) DeleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	delete(s.tasks, id)
	return nil
}

func main() {
	var store Storage = NewMemoryStogare()
	created, err := store.CreateTask(Task{
		Title:       "Сделать ИДЗ по матану",
		Description: "разобраться с методами решения",
		Priority:    1,
		Deadline:    time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		fmt.Printf("Ошибка создания: %v\n", err)
		return
	}
	fmt.Printf("Создана задача с ID: %d\n", created.ID)

	task, err := store.GetTask(created.ID)
	if err != nil {
		fmt.Printf("Ошибка получения: %v\n", err)
	}
	fmt.Println("Получена задача: %v\n", task)

	task.Completed = true
	if err := store.UpdateTask(task); err != nil {
		fmt.Printf("Ошибка обновления: %v\n", err)
		return
	}
	fmt.Println("Задача отмечена как выполненная.")

	if err := store.DeleteTask(task.ID); err != nil {
		fmt.Printf("Ошибка удаления: %v\n", err)
		return
	}
	fmt.Println("Задача успешно удалена.")

	_, err = store.GetTask(task.ID)
	if errors.Is(err, ErrTaskNotFound) {
		fmt.Println("Проверка пройдена: задача действительно отсутствует в хранилище.")
	}
}
