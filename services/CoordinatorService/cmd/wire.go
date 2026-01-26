//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ReilEgor/Vaca/services/CoordinatorService/config"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest"
	handler "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest/handlers"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/usecase"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/repository/redis"
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
var InfraSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedisClient,
	redis.NewRedisTokenRepository,
)

type App struct {
	Logic      domain.CoordinatorUsecase
	Server     *rest.GinServer
	Repository domain.StatusRepository
}

func InitializeApp() (*App, func(), error) {
	wire.Build(
		InfraSet,
		UsecaseSet,
		RestSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
