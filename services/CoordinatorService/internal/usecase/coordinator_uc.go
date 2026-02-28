package usecase

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/google/uuid"
)

type CoordinatorInteractor struct {
	logger       *slog.Logger
	stateClient  domain.StatusRepository
	broker       domain.TaskPublisher
	searchClient domain.SearchRepository
}

const (
	componentCoordinatorUC = "coordinator_uc"
	routingKeyPrefix       = "scraper."
)

func NewCoordinatorUsecase(
	stateClient domain.StatusRepository,
	br domain.TaskPublisher,
	searcher domain.SearchRepository,
) *CoordinatorInteractor {
	return &CoordinatorInteractor{
		stateClient:  stateClient,
		broker:       br,
		searchClient: searcher,
		logger:       slog.With(slog.String("component", componentCoordinatorUC)),
	}
}

func (uc *CoordinatorInteractor) GetTaskStatus(ctx context.Context, taskID string) (*outPkg.Task, error) {
	parsedID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidTaskID, err)
	}

	ans, err := uc.stateClient.Get(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task status: %w", err)
	}

	status, ok := ans["status"]
	if !ok {
		uc.logger.Warn("status key missing in state", slog.String("task_id", taskID))
		return nil, domain.ErrTaskNotFound
	}

	return &outPkg.Task{
		ID:        parsedID,
		Status:    status,
		CreatedAt: time.Time{},
	}, nil
}

func (uc *CoordinatorInteractor) CreateTask(ctx context.Context, keywords []string, sources []string) (*outPkg.Task, error) {
	searchKey := generateSearchKey(keywords, sources)
	existingID, err := uc.stateClient.GetIDByHash(ctx, searchKey)
	if err == nil && existingID != "" {
		id, err := uuid.Parse(existingID)
		if err != nil {
			return nil, fmt.Errorf("%w: existing task: %v", domain.ErrInvalidTaskID, err)
		}
		return &outPkg.Task{ID: id, Status: "exists"}, nil
	}

	taskID := uuid.New()
	now := time.Now()

	if err := uc.stateClient.Set(ctx, taskID.String(), searchKey, len(sources)); err != nil {
		return nil, fmt.Errorf("%w: state set: %v", domain.ErrFailedToCreateTask, err)
	}

	var publishErrors []string
	for _, source := range sources {
		msg := outPkg.ScrapeTask{
			ID:      taskID,
			Keyword: keywords,
			Source:  source,
		}

		if err := uc.broker.PublishTask(ctx, msg, routingKeyPrefix+source); err != nil {
			uc.logger.Error("failed to publish task",
				slog.String("task_id", taskID.String()),
				slog.String("source", source),
				slog.Any("error", err),
			)
			publishErrors = append(publishErrors, source)
		}
	}
	if len(publishErrors) > 0 {
		uc.logger.Warn("failed to publish to some sources",
			slog.String("task_id", taskID.String()),
			slog.Any("failed_sources", publishErrors),
		)
	}
	return &outPkg.Task{
		ID:        taskID,
		Status:    "created",
		CreatedAt: now,
	}, nil
}

func (uc *CoordinatorInteractor) GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, int64, error) {
	vacancies, err := uc.searchClient.GetVacancies(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", domain.ErrSearchFailed, err)
	}
	return vacancies, int64(len(vacancies)), nil
}

func (uc *CoordinatorInteractor) GetAvailableSources(ctx context.Context) ([]outPkg.Source, int64, error) {
	sources, err := uc.stateClient.GetSources(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get available sources: %v", err)
	}
	return sources, int64(len(sources)), err
}

func generateSearchKey(keywords, sources []string) string {
	sort.Strings(keywords)
	sort.Strings(sources)
	input := strings.Join(keywords, ",") + "|" + strings.Join(sources, ",")
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
