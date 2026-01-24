package usecase

import (
	"context"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/google/uuid"
)

type CoordinatorInteractor struct {
	//TODO: Add dependencies
}

func NewCoordinatorUsecase() *CoordinatorInteractor {
	//TODO: Implement constructor
	return &CoordinatorInteractor{}
}

func (uc *CoordinatorInteractor) GetTaskStatus(ctx context.Context, taskID string) (*outPkg.Task, error) {
	return nil, nil
}
func (uc *CoordinatorInteractor) CreateTask(ctx context.Context, keywords []string, sources []string) (uuid.UUID, error) {
	return uuid.New(), nil
}
func (uc *CoordinatorInteractor) GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, int64, error) {
	return nil, 0, nil
}
func (uc *CoordinatorInteractor) GetAvailableSources(ctx context.Context) ([]*outPkg.Source, int64, error) {
	return nil, 0, nil
}
