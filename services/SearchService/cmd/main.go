package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/ReilEgor/Vaca/services/SearchService/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	logger = slog.With(slog.String("service", "main"))

	app, cleanup, err := InitializeApp(
		config.ElasticSearchURL(config.GetEnv("ELASTICSEARCH_URL", "http://elasticsearch:9200")),
		config.GRPCPort(config.GetEnv("GRPC_PORT", "50051")),
	)
	if cleanup != nil {
		defer cleanup()
	}

	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	if err := app.Server.Start(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
