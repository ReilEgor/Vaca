package handler

import (
	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Keywords []string `json:"keywords" binding:"required,gt=0"`
	Sources  []string `json:"sources" binding:"required,gt=0"`
}

type CreateTaskResponse struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) GetTaskStatus(c *gin.Context) {

}
func (h *Handler) CreateTask(c *gin.Context) {

}
