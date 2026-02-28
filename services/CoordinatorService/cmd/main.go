package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ReilEgor/Vaca/pkg"
	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	_ "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest/handlers"
	"github.com/caarlos0/env/v11"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	taskQueue := config.PublisherQueueName(outPkg.RabbitMQTaskQueue)

	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		logger.Error("failed to parse environment variables", slog.Any("error", err))
		os.Exit(1)
	}

	app, cleanup, err := InitializeApp(config.RabbitURL(cfg.RabbitURL), config.SearchClientAddr(cfg.SearchServiceAddress), taskQueue, config.StateClientAddr(cfg.StateServiceAddress))
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		logger.Info("running cleanup...")
		cleanup()
	}()

	serverErr := make(chan error, 1)
	go func() {
		if err := app.Server.Run(":" + cfg.RestApiPort); err != nil {
			serverErr <- fmt.Errorf("server run: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-serverErr:
		logger.Error("server error", slog.Any("error", err))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.Server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("server stopped gracefully")
}
