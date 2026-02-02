package main

import (
	"context"
	_ "github.com/ReilEgor/Vaca/pkg"
	outPkg "github.com/ReilEgor/Vaca/pkg"
	rabbitmq "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/broker/rabbitmq"
	elastic "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/repository/elasticsearch"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	taskQueue := rabbitmq.PublisherQueueName(outPkg.RabbitMQTaskQueue)
	rabbitURL := os.Getenv("RABBIT_URL")
	elasticURL := os.Getenv("ELASTICSEARCH_URL")

	app, cleanup, err := InitializeApp(rabbitmq.RabbitURL(rabbitURL), elastic.ElasticSearchURL(elasticURL), taskQueue)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}
	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		port := os.Getenv("HTTP_PORT")
		if port == "" {
			port = "8080"
		}
		if err := app.Server.Run(":" + port); err != nil {
			logger.Error("failed to start server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")
}
