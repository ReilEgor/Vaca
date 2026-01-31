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
type ResultPublisher interface {
	Publish(result []byte) error
}

//go:generate mockery --name DataProcessorUsecase --output ../mocks/domain --outpkg domain --case=underscore
type DataProcessorUsecase interface {
	Process(ctx context.Context, vacancies outPkg.ScrapeResult) error
}

//go:generate mockery --name TaskCache --output ../mocks/domain --outpkg domain --case=underscore
type TaskCache interface {
	IncrementCompleted(ctx context.Context, taskID string) (int64, error)
	SetStatus(ctx context.Context, taskID string, status string) error
	GetTotal(ctx context.Context, taskID string) (int64, error)
}

//go:generate mockery --name VacancyRepository --output ../mocks/domain --outpkg domain --case=underscore
type VacancyRepository interface {
	SaveBatch(ctx context.Context, result outPkg.ScrapeResult) error
}

//go:generate mockery --name VacancySearchRepository --output ../mocks/domain --outpkg domain --case=underscore
type VacancySearchRepository interface {
	Index(ctx context.Context, v outPkg.Vacancy, taskID string) error
	IndexBatch(ctx context.Context, vacancies outPkg.ScrapeResult) error
	Search(ctx context.Context, query string) ([]outPkg.Vacancy, error)
}
