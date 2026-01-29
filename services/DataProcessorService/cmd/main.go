package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ReilEgor/Vaca/pkg"
	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/broker/rabbitmq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	logger = slog.With(slog.String("service", "data-processor-service"))
	dsn := os.Getenv("DB_SOURCE")
	rabbitURL := os.Getenv("RABBIT_URL")
	app, cleanup, err := InitializeApp(dsn, rabbitmq.RabbitURL(rabbitURL), outPkg.RabbitMQVacancyQueue, logger)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}

	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("monitor service is starting")
	go func() {
		logger.Info("starting consumer")
		app.Subscriber.Listen(ctx)
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")
}
