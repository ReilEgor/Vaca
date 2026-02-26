package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ReilEgor/Vaca/services/NotificationService/internal/config"
	"github.com/caarlos0/env/v11"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	logger = slog.With(slog.String("service", "main"))
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		logger.Error("failed to parse config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	app, cleanup, err := InitializeApp(ctx, config.RabbitMQURL(cfg.RabbitURL), "notification_queue")
	if err != nil {
		logger.Error("failed to initialize app",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer cleanup()

	logger.Debug("notification service is starting")
	go func() {
		err = app.listener.Listen(ctx)
		if err != nil {
			logger.Error("subscriber stopped with error",
				slog.Any("error", err),
			)
			return
		}
	}()
	<-ctx.Done()
	logger.Debug("shutting down gracefully")
}
