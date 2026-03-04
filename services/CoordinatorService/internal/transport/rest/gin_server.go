package rest

import (
	"log/slog"

	handler "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest/handlers"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/usecase"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	router *gin.Engine
	uc     *usecase.CoordinatorInteractor
	logger *slog.Logger
}

func NewGinServer(uc *usecase.CoordinatorInteractor) *GinServer {
	router := gin.New()
	logger := slog.With(slog.String("component", "gin_server"))

	SetupMiddleware(router, logger)

	s := &GinServer{
		router: router,
		uc:     uc,
		logger: logger,
	}

	h := handler.NewHandler(uc)
	h.InitRoutes(s.router)

	return s
}

func (s *GinServer) Run(port string) error {
	return s.router.Run(port)
}
