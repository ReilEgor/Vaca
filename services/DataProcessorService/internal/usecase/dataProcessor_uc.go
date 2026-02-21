package usecase

import (
	"context"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
)

type DataProcessorInteractor struct {
	logger           *slog.Logger
	stateRepository  domain.StateRepository
	repository       domain.VacancyRepository
	searchRepository domain.SearchRepository
	// publisher *rabbitmq.Publisher
}

func NewDataProcessorInteractor(stateClient domain.StateRepository, repository domain.VacancyRepository, searchRepository domain.SearchRepository) *DataProcessorInteractor {
	return &DataProcessorInteractor{
		logger:           slog.With(slog.String("component", "dataProcessorInteractor")),
		stateRepository:  stateClient,
		repository:       repository,
		searchRepository: searchRepository,
	}
}

func (i *DataProcessorInteractor) Process(ctx context.Context, vacancies outPkg.ScrapeResult) error {
	taskID := vacancies.TaskID
	current, err := i.stateRepository.IncrementCompleted(ctx, taskID.String())
	if err != nil {
		return err
	}
	total, err := i.stateRepository.GetTotal(ctx, taskID.String())
	if err != nil {
		return err
	}

	err = i.repository.SaveBatch(ctx, vacancies)
	if err != nil {
		return err
	}

	if current >= total {
		i.logger.Debug("all vacancies processed", slog.String("task_id", taskID.String()))
		err := i.stateRepository.SetStatus(ctx, taskID.String(), "completed")
		if err != nil {
			return err
		}
		err = i.searchRepository.SetVacancies(ctx, vacancies)
		if err != nil {
			return err
		}
	}
	return nil
}
