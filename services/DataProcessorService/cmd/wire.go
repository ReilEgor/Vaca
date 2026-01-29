//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"

	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/config"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/repository/postgres"
	redis "github.com/ReilEgor/Vaca/services/DataProcessorService/internal/repository/redis"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/usecase"
	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	usecase.NewDataProcessorInteractor,
	wire.Bind(new(domain.DataProcessorUsecase), new(*usecase.DataProcessorInteractor)),
)

var BrokerSet = wire.NewSet(
	rabbitmq.NewRabbitMQConn,
	rabbitmq.NewRabbitMQChannel,
	rabbitmq.NewTaskSubscriber,

	wire.Bind(new(domain.DataSubscriber), new(*rabbitmq.DataSubscriber)),
)
var RepositorySet = wire.NewSet(
	postgres.NewPostgresDB,
	postgres.NewVacancyRepository,
)
var InfraSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedisClient,
	redis.NewRedisTokenRepository,
)

type App struct {
	Logic      domain.DataProcessorUsecase
	Repository domain.VacancyRepository
	Subscriber domain.DataSubscriber
	Cache      domain.TaskCache
	//Broker     domain.TaskPublisher
}

func InitializeApp(
	dsn string,
	rabbitURL rabbitmq.RabbitURL,
	qName rabbitmq.SubscriberQueueName,
	logger *slog.Logger,
) (*App, func(), error) {
	wire.Build(
		InfraSet,
		RepositorySet,
		UsecaseSet,
		BrokerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
