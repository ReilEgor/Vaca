package domain

import (
	"context"

	outPkg "github.com/ReilEgor/Vaca/pkg"
)

//go:generate mockery --name DataSubscriber --output ../mocks/domain --outpkg domain --case=underscore
type DataSubscriber interface {
	Listen(ctx context.Context) error
}

//go:generate mockery --name ResultPublisher --output ../mocks/domain --outpkg domain --case=underscore
type DataPublisher interface {
	Publish(ctx context.Context, message outPkg.ScrapeResult) error
}

//go:generate mockery --name DataProcessorUsecase --output ../mocks/domain --outpkg domain --case=underscore
type DataProcessorUsecase interface {
	Process(ctx context.Context, vacancies outPkg.ScrapeResult) error
}

//go:generate mockery --name VacancyRepository --output ../mocks/domain --outpkg domain --case=underscore
type VacancyRepository interface {
	SaveBatch(ctx context.Context, result outPkg.ScrapeResult) error
}

//go:generate mockery --name StatusRepository --output ../mocks/domain --outpkg domain --case=underscore
type StateRepository interface {
	IncrementCompleted(ctx context.Context, taskID string) (int64, error)
	SetStatus(ctx context.Context, taskID string, status string) error
	GetTotal(ctx context.Context, taskID string) (int64, error)
}

//go:generate mockery --name SearchRepository --output ../mocks/domain --outpkg domain --case=underscore
type SearchRepository interface {
	SetVacancies(ctx context.Context, vacancies outPkg.ScrapeResult) error
	GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, error)
}
