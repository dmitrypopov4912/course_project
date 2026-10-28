package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"task-manager/internal/app"
	"task-manager/internal/presentations/task_presentation"
	task_services "task-manager/internal/services"
	"task-manager/internal/storage/postgres_storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err.Error())
	}
	defer pool.Close()
	pgStorage := postgres_storage.NewStorage(pool)
	taskService := task_services.NewService(pgStorage)
	taskPresentation := task_presentation.NewPresentation(taskService)
	e := app.BuildApp(taskPresentation)
	defer e.Close()
	if err := e.Start(":8080"); err != nil {
		panic(fmt.Errorf("cant start server: %w", err))
	}
}
