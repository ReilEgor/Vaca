package usecase

import (
	"context"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
)

type NotificationUsecase struct {
	logger *slog.Logger
}

func NewNotificationUsecase() *NotificationUsecase {
	return &NotificationUsecase{
		logger: slog.With(slog.String("component", "dataProcessorInteractor")),
	}
}

func (u *NotificationUsecase) SendNotification(ctx context.Context, message outPkg.ScrapeResult) error {
	u.logger.Info("NotificationUsecase is listening for messages...")
	return nil
}
