//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"

	"github.com/ReilEgor/Vaca/services/DouScraper/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/config"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/domain"
	redis "github.com/ReilEgor/Vaca/services/DouScraper/internal/repository/redis"
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
var InfraSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedisClient,
	redis.NewRedisScraperRepo,
	wire.Bind(new(domain.SourceRepository), new(*redis.RedisScraperRepo)),
)

type App struct {
	Logic      domain.ScraperUsecase
	Subscriber domain.TaskSubscriber
	Repository domain.SourceRepository
}

func InitializeApp(
	rabbitURL rabbitmq.RabbitURL,
	qName rabbitmq.SubscriberQueueName,
	rKey rabbitmq.SubscriberRoutingKey,
	exch rabbitmq.SubscriberExchange,
	publishQName rabbitmq.PublisherQueueName,
	logger *slog.Logger,
) (*App, func(), error) {
	wire.Build(
		InfraSet,
		BrokerSet,
		UsecaseSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
