package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ReilEgor/Vaca/pkg"
	outPkg "github.com/ReilEgor/Vaca/pkg"
	rabbitmq "github.com/ReilEgor/Vaca/services/DouScraper/internal/broker/rabbitmq"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rabbitURL := os.Getenv("RABBIT_URL")
	app, cleanup, err := InitializeApp(
		rabbitmq.RabbitURL(rabbitURL),
		rabbitmq.SubscriberQueueName("dou_tasks"),
		rabbitmq.SubscriberRoutingKey("scraper.dou.ua"),
		rabbitmq.SubscriberExchange(outPkg.RabbitMQExchangeName),
		logger,
	)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}

	defer cleanup()

	if err := app.Subscriber.Listen(ctx); err != nil {
		logger.Error("subscriber stop", slog.Any("error", err))
	}

	<-ctx.Done()
	logger.Info("shutting down gracefully")
}
