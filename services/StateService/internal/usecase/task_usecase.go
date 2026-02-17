package usecase

import (
	"context"
	"log/slog"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/StateService/internal/domain"
)

type TaskInteractor struct {
	logger *slog.Logger
	repo   domain.StateRepository
}

func NewTaskUsecase(repo domain.StateRepository) *TaskInteractor {
	return &TaskInteractor{
		logger: slog.With(slog.String("component", "task_uc")),
		repo:   repo,
	}
}

func (u *TaskInteractor) Set(ctx context.Context, taskID string, searchKey string, totalSources int, ttl time.Duration) error {
	err := u.repo.Set(ctx, taskID, searchKey, totalSources, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (u *TaskInteractor) Get(ctx context.Context, taskID string) map[string]string {
	task := u.repo.Get(ctx, taskID)
	return task
}

func (u *TaskInteractor) GetIDByHash(ctx context.Context, searchKey string) (string, error) {
	id, err := u.repo.GetIDByHash(ctx, searchKey)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (u *TaskInteractor) GetSources(ctx context.Context) ([]outPkg.Source, error) {
	sources, err := u.repo.GetSources(ctx)
	if err != nil {
		return nil, err
	}
	return sources, nil
}

func (u *TaskInteractor) SetStatus(ctx context.Context, taskID string, status string) error {
	err := u.repo.SetStatus(ctx, taskID, status)
	if err != nil {
		return err
	}
	return nil
}

func (u *TaskInteractor) GetTotal(ctx context.Context, taskID string) (int64, error) {
	total, err := u.repo.GetTotal(ctx, taskID)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (u *TaskInteractor) IncrementCompleted(ctx context.Context, taskID string) (int64, error) {
	total, err := u.repo.IncrementCompleted(ctx, taskID)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (u *TaskInteractor) Register(ctx context.Context, source outPkg.Source) error {
	err := u.repo.Register(ctx, source)
	if err != nil {
		return err
	}
	return nil
}
