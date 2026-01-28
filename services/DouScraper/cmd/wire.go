//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"

	"github.com/ReilEgor/Vaca/services/DouScraper/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/domain"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/usecase"
	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	usecase.NewDouInteractor,
	wire.Bind(new(domain.ScraperUsecase), new(*usecase.DouInteractor)),
)

var BrokerSet = wire.NewSet(
	rabbitmq.NewRabbitMQConn,
	rabbitmq.NewRabbitMQChannel,
	rabbitmq.NewTaskSubscriber,
	wire.Bind(new(domain.TaskSubscriber), new(*rabbitmq.TaskSubscriber)),
)

type App struct {
	Logic      domain.ScraperUsecase
	Subscriber domain.TaskSubscriber
}

func InitializeApp(
	rabbitURL rabbitmq.RabbitURL,
	qName rabbitmq.SubscriberQueueName,
	rKey rabbitmq.SubscriberRoutingKey,
	exch rabbitmq.SubscriberExchange,
	logger *slog.Logger,
) (*App, func(), error) {
	wire.Build(
		BrokerSet,
		UsecaseSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
