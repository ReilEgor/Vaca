package usecase

import (
	"context"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	"golang.org/x/sync/errgroup"
)

type DataProcessorInteractor struct {
	logger           *slog.Logger
	stateRepository  domain.StateRepository
	repository       domain.VacancyRepository
	searchRepository domain.SearchRepository
	dataPublisher    domain.DataPublisher
	// publisher *rabbitmq.Publisher
}

func NewDataProcessorInteractor(stateClient domain.StateRepository, repository domain.VacancyRepository, searchRepository domain.SearchRepository, dataPublisher domain.DataPublisher) *DataProcessorInteractor {
	return &DataProcessorInteractor{
		logger:           slog.With(slog.String("component", "dataProcessorInteractor")),
		stateRepository:  stateClient,
		repository:       repository,
		searchRepository: searchRepository,
		dataPublisher:    dataPublisher,
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
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return i.repository.SaveBatch(gCtx, vacancies)
	})
	g.Go(func() error {
		return i.dataPublisher.Publish(gCtx, vacancies)
	})
	if err := g.Wait(); err != nil {
		return fmt.Errorf("batch processing failed: %w", err)
	}
	if current >= total {
		g, gCtx := errgroup.WithContext(ctx)
		i.logger.Debug("all vacancies processed", slog.String("task_id", taskID.String()))
		g.Go(func() error {
			return i.searchRepository.SetVacancies(gCtx, vacancies)
		})
		g.Go(func() error {
			return i.stateRepository.SetStatus(gCtx, taskID.String(), "completed")
		})
		if err := g.Wait(); err != nil {
			i.logger.Error("failed to finalize task", slog.Any("error", err))
			return err
		}
	}
	return nil
}
