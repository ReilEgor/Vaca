package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ReilEgor/Vaca/pkg"
	outPkg "github.com/ReilEgor/Vaca/pkg"
	rabbitmq "github.com/ReilEgor/Vaca/services/DouScraper/internal/broker/rabbitmq"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	logger = slog.With(slog.String("service", "main"))

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
		rabbitmq.PublisherQueueName(outPkg.RabbitMQVacancyQueue),
		logger,
	)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}

	source := outPkg.Source{
		ID:   uuid.New(),
		Name: "dou.ua",
		URL:  "https://dou.ua",
	}
	err = app.Repository.Register(ctx, source, time.Hour*24)
	if err != nil {
		logger.Error("failed to register source", slog.Any("error", err))
		return
	}

	defer cleanup()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8080", nil); err != nil {
			slog.Error("prometheus server failed", "error", err)
		}
	}()

	if err := app.Subscriber.Listen(ctx); err != nil {
		logger.Error("subscriber stop", slog.Any("error", err))
	}

	<-ctx.Done()
	logger.Debug("shutting down gracefully")
}
