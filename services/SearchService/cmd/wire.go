//go:build wireinject
// +build wireinject

package main

import (
	elasticGen "github.com/ReilEgor/Vaca/services/SearchService/api/proto/gen/elastic"
	"github.com/ReilEgor/Vaca/services/SearchService/internal/config"
	"github.com/ReilEgor/Vaca/services/SearchService/internal/domain"
	elastic "github.com/ReilEgor/Vaca/services/SearchService/internal/repository/elasticsearch"
	grpcH "github.com/ReilEgor/Vaca/services/SearchService/internal/transport/grpc"
	"github.com/ReilEgor/Vaca/services/SearchService/internal/usecase"
	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	usecase.NewSearchUsecase,
	wire.Bind(new(domain.SearchUsecase), new(*usecase.SearchInteractor)),
)

var GRPCSet = wire.NewSet(
	grpcH.NewSearchHandler,
	grpcH.NewServer,
	wire.Bind(new(elasticGen.SearchServiceServer), new(*grpcH.SearchHandler)),
)

var ElasticSet = wire.NewSet(
	elastic.NewElasticClient,
	elastic.NewElasticRepository,
)

type App struct {
	Repository domain.SearchRepository
	Server     *grpcH.Server
}

func InitializeApp(
	searchRepoURL config.ElasticSearchURL,
	grpcPort config.GRPCPort,
) (*App, func(), error) {
	wire.Build(
		ElasticSet,
		UsecaseSet,
		GRPCSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
