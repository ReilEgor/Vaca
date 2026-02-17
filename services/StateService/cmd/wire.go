//go:build wireinject
// +build wireinject

package main

import (
	taskstate "github.com/ReilEgor/Vaca/services/StateService/api/proto/gen/task_state"
	"github.com/ReilEgor/Vaca/services/StateService/internal/config"
	"github.com/ReilEgor/Vaca/services/StateService/internal/domain"
	"github.com/ReilEgor/Vaca/services/StateService/internal/repository/redis"
	grpcH "github.com/ReilEgor/Vaca/services/StateService/internal/transport/grpc"
	"github.com/ReilEgor/Vaca/services/StateService/internal/usecase"
	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	usecase.NewTaskUsecase,
	wire.Bind(new(domain.StateUsecase), new(*usecase.TaskInteractor)),
)

var GRPCSet = wire.NewSet(
	grpcH.NewTaskStateHandler,
	grpcH.NewServer,
	wire.Bind(new(taskstate.TaskStateServiceServer), new(*grpcH.TaskStateHandler)),
)

var RedisSet = wire.NewSet(
	redis.NewRedisClient,
	redis.NewRedisTokenRepository,
)

type App struct {
	Repository domain.StateRepository
	Server     *grpcH.Server
}

func InitializeApp(
	host config.Host,
	redisPort config.RedisPort,
	password config.Password,
	db config.DB,
	grpcPort config.GRPCPort,
) (*App, func(), error) {
	wire.Build(
		RedisSet,
		GRPCSet,
		UsecaseSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
