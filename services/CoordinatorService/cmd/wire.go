//go:build wireinject
// +build wireinject

package main

import (
	rabbitmq "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	redis2 "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/repository/redis"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest"
	handler "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest/handlers"
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

var InfraSet = wire.NewSet(
	config.NewConfig,
	redis2.NewRedisClient,
	redis2.NewRedisTokenRepository,
)

type App struct {
	Logic      domain.CoordinatorUsecase
	Server     *rest.GinServer
	Repository domain.StatusRepository
	//Broker     domain.TaskPublisher
}

func InitializeApp(rabbitURL rabbitmq.RabbitURL, taskQueue rabbitmq.PublisherQueueName) (*App, func(), error) {
	wire.Build(
		InfraSet,
		BrokerSet,
		UsecaseSet,
		RestSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
