package domain

import (
	"context"

	outPkg "github.com/ReilEgor/Vaca/pkg"
)

type NotificationUsecase interface {
	// TODO: create struct for message with more fields (e.g. type, recipient, etc.)
	SendNotification(ctx context.Context, message outPkg.ScrapeResult) error
}

type NotificationSubscriber interface {
	Listen(ctx context.Context) error
}
