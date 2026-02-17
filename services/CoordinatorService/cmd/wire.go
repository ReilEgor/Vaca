//go:build wireinject
// +build wireinject

package main

import (
	rabbitmq "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	elastic "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/repository/elasticsearch"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest"
	handler "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest/handlers"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/stateClient"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/usecase"
	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	usecase.NewCoordinatorUsecase,
	wire.Bind(new(domain.CoordinatorUsecase), new(*usecase.CoordinatorInteractor)),
)

var RestSet = wire.NewSet(
	rest.NewGinServer,
	handler.NewHandler,
)

var BrokerSet = wire.NewSet(
	rabbitmq.NewRabbitMQConn,
	rabbitmq.NewRabbitMQChannel,
	rabbitmq.NewPublisher,
	wire.Bind(new(domain.TaskPublisher), new(*rabbitmq.Publisher)),
)

var ElasticSet = wire.NewSet(
	elastic.NewElasticClient,
	elastic.NewElasticRepository,
	wire.Bind(new(domain.VacancySearchRepository), new(*elastic.ElasticRepository)),
)

var StateClientSet = wire.NewSet(
	stateClient.NewStateClient,
)

type App struct {
	Logic      domain.CoordinatorUsecase
	Server     *rest.GinServer
	SearchRepo domain.VacancySearchRepository
}

func InitializeApp(
	rabbitURL rabbitmq.RabbitURL,
	searchRepoURL elastic.ElasticSearchURL,
	taskQueue rabbitmq.PublisherQueueName,
	stateClientAddr config.StateClientAddr,
) (*App, func(), error) {
	wire.Build(
		StateClientSet,
		BrokerSet,
		ElasticSet,
		UsecaseSet,
		RestSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
