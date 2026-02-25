package domain

import "context"

type NotificationUsecase interface {
	// TODO: create struct for message with more fields (e.g. type, recipient, etc.)
	SendNotification(ctx context.Context, message string) error
}

type NotificationSubscriber interface {
	Listen(ctx context.Context) error
}
