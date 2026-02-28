package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ReilEgor/Vaca/pkg"
	outPkg "github.com/ReilEgor/Vaca/pkg"
	rabbitmq "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/broker/rabbitmq"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	_ "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/transport/rest/handlers"
	"github.com/joho/godotenv"
)

// @title           Coordinator Service API
// @version         1.0
// @description     This is the Coordinator Service for the Vaca project.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Egor Reil
// @contact.url    https://github.com/ReilEgor

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http https

// @accept    json
// @produce   json
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	taskQueue := rabbitmq.PublisherQueueName(outPkg.RabbitMQTaskQueue)
	rabbitURL := os.Getenv("RABBIT_URL")
	searchServiceAddress := os.Getenv("SEARCH_SERVICE_ADDRESS")
	stateServiceAddress := os.Getenv("STATE_SERVICE_ADDRESS")
	app, cleanup, err := InitializeApp(rabbitmq.RabbitURL(rabbitURL), config.SearchClientAddr(searchServiceAddress), taskQueue, config.StateClientAddr(stateServiceAddress))
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
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")
}
