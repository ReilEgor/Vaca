package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ReilEgor/Vaca/pkg"
	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	logger = slog.With(slog.String("service", "main"))

	dsn := os.Getenv("DB_SOURCE")
	rabbitURL := os.Getenv("RABBIT_URL")

	stateServiceAddress := os.Getenv("STATE_SERVICE_ADDRESS")
	searchServiceAddress := os.Getenv("SEARCH_SERVICE_ADDRESS")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, cleanup, err := InitializeApp(ctx, dsn, rabbitmq.RabbitURL(rabbitURL), outPkg.RabbitMQVacancyQueue, logger, config.StateClientAddr(stateServiceAddress), config.SearchClientAddr(searchServiceAddress))
	if err != nil {
		logger.Error("failed to initialize app",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	defer cleanup()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8080", nil); err != nil {
			slog.Error("prometheus server failed", "error", err)
		}
	}()

	logger.Debug("data-processor service is starting")
	go func() {
		logger.Debug("starting subscriber")
		err := app.Subscriber.Listen(ctx)
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
