package domain

import (
	"context"
	outPkg "github.com/ReilEgor/Vaca/pkg"
)

//go:generate mockery --name CoordinatorUsecase --output ../mocks/domain --outpkg domain --case=underscore
type CoordinatorUsecase interface {
	GetTaskStatus(ctx context.Context, taskID string) (*outPkg.Task, error)
	CreateTask(ctx context.Context, keywords []string, sources []outPkg.Source) (string, error)
	GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, int64, error)
	GetAvailableSources(ctx context.Context) ([]*outPkg.Source, int64, error)
}
