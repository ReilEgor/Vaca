//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/ReilEgor/Vaca/services/NotificationService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/NotificationService/internal/config"
	"github.com/ReilEgor/Vaca/services/NotificationService/internal/domain"
	"github.com/ReilEgor/Vaca/services/NotificationService/internal/usecase"
	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	usecase.NewNotificationUsecase,
	wire.Bind(new(domain.NotificationUsecase), new(*usecase.NotificationUsecase)),
)

var ListenerSet = wire.NewSet(
	rabbitmq.NewRabbitMQConn,
	rabbitmq.NewRabbitMQChannel,
	rabbitmq.NewNotificationSubscriber,
	wire.Bind(new(domain.NotificationSubscriber), new(*rabbitmq.NotificationSubscriber)),
)

type App struct {
	listener domain.NotificationSubscriber
}

func InitializeApp(
	ctx context.Context,
	rabbitURL config.RabbitMQURL,
	qName config.SubscriberQueueName,
) (*App, func(), error) {
	wire.Build(
		UsecaseSet,
		ListenerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
