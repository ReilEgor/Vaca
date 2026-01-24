package pkg

import "context"

type Vacancy struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Company     string `json:"company"`
	Salary      string `json:"salary"`
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
