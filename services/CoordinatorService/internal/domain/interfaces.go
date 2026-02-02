package domain

import (
	"context"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/google/uuid"
)

//go:generate mockery --name CoordinatorUsecase --output ../mocks/domain --outpkg domain --case=underscore
type CoordinatorUsecase interface {
	GetTaskStatus(ctx context.Context, taskID string) (*outPkg.Task, error)
	CreateTask(ctx context.Context, keywords []string, sources []string) (*uuid.UUID, error)
	GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, int64, error)
	GetAvailableSources(ctx context.Context) ([]outPkg.Source, int64, error)
}

//go:generate mockery --name StatusRepository --output ../mocks/domain --outpkg domain --case=underscore
type StatusRepository interface {
	Set(ctx context.Context, taskID string, searchKey string, totalSources int, ttl time.Duration) error
	Get(ctx context.Context, taskID string) map[string]string
	GetIDByHash(ctx context.Context, searchKey string) (string, error)
	GetSources(ctx context.Context) ([]outPkg.Source, error)
}

//go:generate mockery --name TaskPublisher --output ../mocks/domain --outpkg domain --case=underscore
type TaskPublisher interface {
	PublishTask(ctx context.Context, taskMessage outPkg.ScrapeTask, routingKey string) error
}

//go:generate mockery --name VacancySearchRepository --output ../mocks/domain --outpkg domain --case=underscore
type VacancySearchRepository interface {
	Search(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, error)
}
