package handler

import (
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
	TaskID    uuid.UUID `json:"task_id" binding:"required"`
	Status    string    `json:"status" binding:"required"`
	CreatedAt string    `json:"created_at" binding:"required"`
}

func (h *Handler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": outPkg.ErrTaskIDRequired.Error()})
		return
	}

	ctx := c.Request.Context()

	task, err := h.uc.GetTaskStatus(ctx, taskID)
	if err != nil {
		h.logger.Error("failed to get task status", slog.String("id", taskID), slog.Any("error", err))
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrTaskNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, CreateTaskResponse{
		TaskID:    task.ID,
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

	ctx := c.Request.Context()

	taskID, err := h.uc.CreateTask(ctx, req.Keywords, req.Sources)
	if err != nil {
		h.logger.Error("failed to create task", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrFailedToCreateTask.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateTaskResponse{
		TaskID:    *taskID,
		Status:    "created",
		CreatedAt: time.Now().Format(time.RFC3339),
	})
}
