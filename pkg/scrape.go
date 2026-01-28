package pkg

import "github.com/google/uuid"

type ScrapeResult struct {
	TaskID    uuid.UUID
	Vacancies []Vacancy
}
