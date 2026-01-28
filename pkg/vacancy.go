package pkg

import (
	"context"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Link         string    `json:"link"`
	Company      string    `json:"company"`
	Salary       string    `json:"salary"`
	City         string    `json:"city"`
	Requirements string    `json:"requirements"`
	About        string    `json:"about"`
}
type VacancyFilter struct {
	Query     string
	City      string
	IsRemote  bool
	SalaryMin int
	SalaryMax int

	Limit  int
	Offset int
}

type VacancySearcher interface {
	Search(ctx context.Context, keywords []string) ([]*Vacancy, error)
}
