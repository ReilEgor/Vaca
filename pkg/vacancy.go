package pkg

import (
	"context"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID           uuid.UUID `json:"task_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Link         string    `json:"link"`
	Company      string    `json:"company"`
	Salary       string    `json:"salary"`
	Location     string    `json:"location"`
	Requirements string    `json:"requirements"`
	About        string    `json:"about"`
}
type VacancyFilter struct {
	Query        string `form:"query"`
	Location     string `form:"location"`
	Requirements string `form:"requirements"`
	Limit        int    `form:"limit"`
	Offset       int    `form:"offset"`
}

type VacancySearcher interface {
	Search(ctx context.Context, keywords []string) ([]*Vacancy, error)
}
