package usecase

import (
	"context"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
)

type DouInteractor struct {
	logger *slog.Logger
}

func NewDouInteractor() *DouInteractor {
	return &DouInteractor{
		logger: slog.With(slog.String("component", "scheduler")),
	}
}

func (i *DouInteractor) Execute(ctx context.Context, task outPkg.ScrapeTask) error {
	i.logger.Info("starting DouInteractor",
		slog.String("url", task.ID.String()),
		slog.Any("task", task))
	return nil
}
