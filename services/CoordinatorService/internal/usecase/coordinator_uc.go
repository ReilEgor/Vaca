package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/google/uuid"
)

type CoordinatorInteractor struct {
	//TODO: Add dependencies
	logger     *slog.Logger
	statusRepo domain.StatusRepository
}

func NewCoordinatorUsecase(sr domain.StatusRepository) *CoordinatorInteractor {
	return &CoordinatorInteractor{
		statusRepo: sr,
		logger:     slog.With(slog.String("component", "coordinator_uc")),
	}
}

func (uc *CoordinatorInteractor) GetTaskStatus(ctx context.Context, taskID string) (*outPkg.Task, error) {
	status, err := uc.statusRepo.Get(ctx, taskID)
	if err != nil {
		uc.logger.Error("failed to get status from repo", slog.Any("error", err))
		return nil, domain.ErrTaskNotFound
	}

	fmt.Printf("status: %s, err: %v\n", status, err)
	parsedID, err := uuid.Parse(taskID)
	if err != nil {
		uc.logger.Error("failed to parse task ID", slog.String("task_id", taskID), slog.Any("error", err))
		return nil, domain.ErrTaskNotFound
	}
	task := &outPkg.Task{
		ID:        parsedID,
		Status:    status,
		CreatedAt: time.Now(),
	}

	return task, nil
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
