package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	handler "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest/handlers"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/usecase"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	router *gin.Engine
	server *http.Server
	uc     *usecase.CoordinatorInteractor
	logger *slog.Logger
}

const (
	componentGinServer = "gin_server"

	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 60 * time.Second
)

func NewGinServer(uc *usecase.CoordinatorInteractor) *GinServer {
	logger := slog.With(slog.String("component", componentGinServer))

	router := gin.New()
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
	s.server = &http.Server{
		Addr:         port,
		Handler:      s.router,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}
	s.logger.Info("starting server", slog.String("port", port))
	return s.router.Run(port)
}

func (s *GinServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return domain.ErrServerNotInitialized
	}
	s.logger.Info("shutting down server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	return nil
}
