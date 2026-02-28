//go:build wireinject
// +build wireinject

package main

import (
	rabbitmq "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest"
	handler "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest/handlers"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/searchClient"
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

var SearchClientSet = wire.NewSet(
	searchClient.NewSearchClient,
	wire.Bind(new(domain.SearchRepository), new(*searchClient.SearchClient)),
)

var StateClientSet = wire.NewSet(
	stateClient.NewStateClient,
	wire.Bind(new(domain.StatusRepository), new(*stateClient.StateClient)),
)

type App struct {
	Logic      domain.CoordinatorUsecase
	Server     *rest.GinServer
	SearchRepo domain.SearchRepository
}

func InitializeApp(
	rabbitURL config.RabbitURL,
	searchClientAddr config.SearchClientAddr,
	taskQueue config.PublisherQueueName,
	stateClientAddr config.StateClientAddr,
) (*App, func(), error) {
	wire.Build(
		StateClientSet,
		BrokerSet,
		SearchClientSet,
		UsecaseSet,
		RestSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
