//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"log/slog"

	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/config"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/repository/postgres"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/transport/searchClient"
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

var SearchClientSet = wire.NewSet(
	searchClient.NewSearchClient,
	wire.Bind(new(domain.SearchRepository), new(*searchClient.SearchClient)),
)

var StateClientSet = wire.NewSet(
	stateClient.NewStateClient,
	wire.Bind(new(domain.StateRepository), new(*stateClient.StateClient)),
)

type App struct {
	Logic      domain.DataProcessorUsecase
	Repository domain.VacancyRepository
	Subscriber domain.DataSubscriber
	SearchRepo domain.SearchRepository
}

func InitializeApp(
	ctx context.Context,
	dsn string,
	rabbitURL rabbitmq.RabbitURL,
	qName rabbitmq.SubscriberQueueName,
	logger *slog.Logger,
	stateClientAddr config.StateClientAddr,
	searchClientAddr config.SearchClientAddr,
) (*App, func(), error) {
	wire.Build(
		SearchClientSet,
		StateClientSet,
		RepositorySet,
		UsecaseSet,
		BrokerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
