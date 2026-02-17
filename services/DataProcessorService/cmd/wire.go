//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"

	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/config"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	elastic "github.com/ReilEgor/Vaca/services/DataProcessorService/internal/repository/elasticsearch"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/repository/postgres"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/transport/stateClient"
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

var ElasticSet = wire.NewSet(
	elastic.NewElasticClient,
	elastic.NewElasticRepository,
)

var StateClient = wire.NewSet(
	stateClient.NewStateClient,
)

type App struct {
	Logic      domain.DataProcessorUsecase
	Repository domain.VacancyRepository
	Subscriber domain.DataSubscriber
	SearchRepo domain.VacancySearchRepository
}

func InitializeApp(
	dsn string,
	rabbitURL rabbitmq.RabbitURL,
	searchRepoURL elastic.ElasticSearchURL,
	qName rabbitmq.SubscriberQueueName,
	logger *slog.Logger,
	stateClientAddr config.StateClientAddr,
) (*App, func(), error) {
	wire.Build(
		StateClient,
		ElasticSet,
		RepositorySet,
		UsecaseSet,
		BrokerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
