package usecase

import (
	"context"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/transport/stateClient"
)

type DataProcessorInteractor struct {
	logger           *slog.Logger
	stateClient      *stateClient.StateClient
	repository       domain.VacancyRepository
	searchRepository domain.VacancySearchRepository
	// publisher *rabbitmq.Publisher
}

func NewDataProcessorInteractor(stateClient *stateClient.StateClient, repository domain.VacancyRepository, searchRepository domain.VacancySearchRepository) *DataProcessorInteractor {
	return &DataProcessorInteractor{
		logger:           slog.With(slog.String("component", "DataProcessorInteractor")),
		stateClient:      stateClient,
		repository:       repository,
		searchRepository: searchRepository,
	}
}

func (i *DataProcessorInteractor) Process(ctx context.Context, vacancies outPkg.ScrapeResult) error {
	taskID := vacancies.TaskID
	current, err := i.stateClient.IncrementCompleted(ctx, taskID.String())
	if err != nil {
		return err
	}
	total, err := i.stateClient.GetTotal(ctx, taskID.String())
	if err != nil {
		return err
	}

	err = i.repository.SaveBatch(ctx, vacancies)
	if err != nil {
		return err
	}

	if current >= total {
		i.logger.Debug("all vacancies processed", slog.String("task_id", taskID.String()))
		err := i.stateClient.SetStatus(ctx, taskID.String(), "completed")
		if err != nil {
			return err
		}
		err = i.searchRepository.IndexBatch(ctx, vacancies)
		if err != nil {
			return err
		}
	}
	return nil
}
