package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Keywords []string `json:"keywords" binding:"required,dive"`
	Sources  []string `json:"sources" binding:"required,dive"`
}

type CreateTaskResponse struct {
	Task      *outPkg.Task `json:"task" binding:"required"`
	Status    string       `json:"status" binding:"required"`
	CreatedAt string       `json:"created_at" binding:"required"`
}

const (
	taskStatusCreated = "created"
)

func (h *Handler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	if _, err := uuid.Parse(taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": outPkg.ErrTaskIDRequired.Error()})
		return
	}

	task, err := h.uc.GetTaskStatus(c.Request.Context(), taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrTaskNotFound.Error()})
			return
		}
		h.logger.Error("failed to get task status",
			slog.String("task_id", taskID),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrFailedToGetTask.Error()})
		return
	}
	c.JSON(http.StatusOK, CreateTaskResponse{
		Task:      task,
		Status:    task.Status,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
	})
}

func (h *Handler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid create task request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": outPkg.ErrInvalidRequest.Error()})
		return
	}

	task, err := h.uc.CreateTask(c.Request.Context(), req.Keywords, req.Sources)
	if err != nil {
		h.logger.Error("failed to create task", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrFailedToCreateTask.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateTaskResponse{
		Task:      task,
		Status:    "created",
		CreatedAt: time.Now().Format(time.RFC3339),
	})
}
