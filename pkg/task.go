package pkg

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID               uuid.UUID `json:"id"`
	KeyWords         []string  `json:"keywords"`
	Status           string    `json:"status"`
	TotalSources     int       `json:"total_sources"`
	CompletedSources int       `json:"completed_sources"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RabbitMQ message structure
type TaskMessage struct {
	TaskID   uuid.UUID
	Keywords []string
	Sources  []string
}
type ScrapeTask struct {
	ID      uuid.UUID
	Keyword []string
	Source  string
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *Task) error
	GetTaskByID(ctx context.Context, id uuid.UUID) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) error
}
