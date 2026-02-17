//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"

	"github.com/ReilEgor/Vaca/services/DouScraper/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/config"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/domain"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/transport/stateClient"
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
	rabbitmq.NewPublisher,
	rabbitmq.NewTaskSubscriber,

	wire.Bind(new(domain.TaskSubscriber), new(*rabbitmq.TaskSubscriber)),
	wire.Bind(new(domain.ResultPublisher), new(*rabbitmq.Publisher)),
)

var StateClientSet = wire.NewSet(
	stateClient.NewStateClient,
)

type App struct {
	Logic      domain.ScraperUsecase
	Subscriber domain.TaskSubscriber
	Repository *stateClient.StateClient
}

func InitializeApp(
	rabbitURL rabbitmq.RabbitURL,
	qName rabbitmq.SubscriberQueueName,
	rKey rabbitmq.SubscriberRoutingKey,
	exch rabbitmq.SubscriberExchange,
	publishQName rabbitmq.PublisherQueueName,
	logger *slog.Logger,
	stateClientAddr config.StateClientAddr,
) (*App, func(), error) {
	wire.Build(
		StateClientSet,
		BrokerSet,
		UsecaseSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
