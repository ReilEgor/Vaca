//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
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

type App struct {
	Logic  domain.CoordinatorUsecase
	Server *rest.GinServer
}

func InitializeApp() (*App, func(), error) {
	wire.Build(
		UsecaseSet,
		RestSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
