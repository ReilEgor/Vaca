package usecase

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"log/slog"
	"sort"
	"strings"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/google/uuid"
)

type CoordinatorInteractor struct {
	//TODO: Add dependencies
	logger     *slog.Logger
	statusRepo domain.StatusRepository
	broker     domain.TaskPublisher
	searcher   domain.VacancySearchRepository
}

func NewCoordinatorUsecase(sr domain.StatusRepository, br domain.TaskPublisher, searcher domain.VacancySearchRepository) *CoordinatorInteractor {
	return &CoordinatorInteractor{
		statusRepo: sr,
		logger:     slog.With(slog.String("component", "coordinator_uc")),
		broker:     br,
		searcher:   searcher,
	}
}

func (uc *CoordinatorInteractor) GetTaskStatus(ctx context.Context, taskID string) (*outPkg.Task, error) {
	ans := uc.statusRepo.Get(ctx, taskID)
	status, ok := ans["status"]
	if !ok {
		uc.logger.Error("status not found in repo", slog.String("task_id", taskID))
		return nil, domain.ErrTaskNotFound
	}

	parsedID, err := uuid.Parse(taskID)
	if err != nil {
		uc.logger.Error("failed to parse task ID", slog.String("task_id", taskID), slog.Any("error", err))
		return nil, domain.ErrTaskNotFound
	}
	task := &outPkg.Task{
		ID:     parsedID,
		Status: status,
		//TODO: Am I needed to store CreatedAt in redis or another storage?
		CreatedAt: time.Time{},
	}

	return task, nil
}
func (uc *CoordinatorInteractor) CreateTask(ctx context.Context, keywords []string, sources []string) (*uuid.UUID, error) {
	searchKey := generateSearchKey(keywords, sources)
	existingID, err := uc.statusRepo.GetIDByHash(ctx, searchKey)
	if err == nil && existingID != "" {
		id, _ := uuid.Parse(existingID)
		return &id, nil
	}

	taskID := uuid.New()
	err = uc.statusRepo.Set(ctx, taskID.String(), searchKey, len(sources), time.Minute*2)
	if err != nil {
		//TODO: create status constants
		uc.logger.Error("failed to set status from repo", slog.Any("error", err))
		//TODO: return proper error
		return nil, domain.ErrTaskNotFound
	}

	for _, source := range sources {
		rKey := "scraper." + source

		msg := outPkg.ScrapeTask{
			ID:      taskID,
			Keyword: keywords,
			Source:  source,
		}

		if err := uc.broker.PublishTask(ctx, msg, rKey); err != nil {
			uc.logger.Error("failed to send task", slog.String("source", source))
		}
	}

	return &taskID, nil
}

func (uc *CoordinatorInteractor) GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, int64, error) {
	vacancies, err := uc.searcher.Search(ctx, filter)
	if err != nil {
		uc.logger.Error("failed to search vacancies", slog.Any("error", err))
		return nil, 0, domain.ErrSearchFailed
	}
	return vacancies, int64(len(vacancies)), nil
}
func (uc *CoordinatorInteractor) GetAvailableSources(ctx context.Context) ([]outPkg.Source, int64, error) {
	//TODO: refactor
	sources, err := uc.statusRepo.GetSources(ctx)
	return sources, int64(len(sources)), err
}

func generateSearchKey(keywords, sources []string) string {
	sort.Strings(keywords)
	sort.Strings(sources)
	input := strings.Join(keywords, ",") + "|" + strings.Join(sources, ",")
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
