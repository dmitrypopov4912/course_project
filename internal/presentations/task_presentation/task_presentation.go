package task_presentation

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"task-manager/internal/dto"
	task_services "task-manager/internal/services"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Presentation struct {
	service *task_services.TaskService
}

func NewPresentation(service *task_services.TaskService) *Presentation {
	return &Presentation{service: service}
}

func (p *Presentation) CreateTask(c echo.Context) error {
	var task dto.TaskRequest
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"err": fmt.Errorf("c.Bind: %w", err).Error()})
	}
	if err := p.service.CreateTask(c.Request().Context(), task); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"err": fmt.Errorf("p.service.CreateTask: %w", err).Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (p *Presentation) GetTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"err": fmt.Errorf("invalid task id: err: %w", err).Error()})
	}
	task, err := p.service.GetTask(c.Request().Context(), int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"err": fmt.Errorf("p.service.GetTask: %w", err).Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"err": fmt.Errorf("p.service.GetTask: %w", err).Error()})
	}
	return c.JSON(http.StatusOK, task)
}

func (p *Presentation) UpdateTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"err": fmt.Errorf("invalid task id: %w", err).Error()})
	}
	var task dto.TaskRequest
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"err": fmt.Errorf("p.service.UpdateTask: %w", err).Error()})
	}
	if err := p.service.UpdateTask(c.Request().Context(), task, int64(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"err": fmt.Errorf("p.service.UpdateTask: %w", err).Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (p *Presentation) DeleteTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"err": fmt.Errorf("invalid task id: err: %w", err).Error()})
	}
	if err := p.service.DeleteTask(c.Request().Context(), int64(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"err": fmt.Errorf("p.service.DeleteTask: %w", err).Error()})
	}
	return c.JSON(http.StatusOK, nil)
}
