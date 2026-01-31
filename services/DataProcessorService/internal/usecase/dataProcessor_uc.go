package usecase

import (
	"context"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
)

type DataProcessorInteractor struct {
	logger           *slog.Logger
	cache            domain.TaskCache
	repository       domain.VacancyRepository
	searchRepository domain.VacancySearchRepository
	//publisher *rabbitmq.Publisher

}

func NewDataProcessorInteractor(cache domain.TaskCache, repository domain.VacancyRepository, searchRepository domain.VacancySearchRepository) *DataProcessorInteractor {
	return &DataProcessorInteractor{
		logger:           slog.With(slog.String("component", "DataProcessorInteractor")),
		cache:            cache,
		repository:       repository,
		searchRepository: searchRepository,
	}
}

func (i *DataProcessorInteractor) Process(ctx context.Context, vacancies outPkg.ScrapeResult) error {
	taskID := vacancies.TaskID
	current, err := i.cache.IncrementCompleted(ctx, taskID.String())
	if err != nil {
		return err
	}
	total, err := i.cache.GetTotal(ctx, taskID.String())
	if err != nil {
		return err
	}

	err = i.repository.SaveBatch(ctx, vacancies)
	if err != nil {
		return err
	}

	if current >= total {
		fmt.Println("All vacancies processed for task:", taskID.String())
		err := i.cache.SetStatus(ctx, taskID.String(), "completed")
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
