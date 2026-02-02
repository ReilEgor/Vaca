package pkg

import (
	"context"

	"github.com/google/uuid"
)

type Source struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	IsActive bool      `json:"is_active"`
}

type SourceRepository interface {
	GetAllActiveSources(ctx context.Context) ([]*Source, error)
}
