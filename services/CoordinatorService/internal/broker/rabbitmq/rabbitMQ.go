package rabbitmq

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ()

const (
	componentConnector = "rabbitmqConnector"
	componentChannel   = "rabbitmqChannel"
)

func NewRabbitMQConn(rabbitURL config.RabbitURL) (*amqp.Connection, func(), error) {
	logger := slog.With(slog.String("component", componentConnector))

	conn, err := amqp.Dial(string(rabbitURL))
	if err != nil {
		logger.Error("failed to connect to rabbitmq",
			slog.Any("error", err),
			slog.String("url", maskRabbitURL(string(rabbitURL))),
		)
		return nil, nil, fmt.Errorf("%w: %v", domain.ErrConnect, err)
	}

	logger.Info("connected to rabbitmq",
		slog.String("url", maskRabbitURL(string(rabbitURL))),
	)

	cleanup := func() {
		logger.Info("closing rabbitmq connection")
		if err := conn.Close(); err != nil {
			logger.Error("failed to close rabbitmq connection",
				slog.Any("error", fmt.Errorf("%w: %w", domain.ErrCloseConn, err)),
			)
		}
	}

	return conn, cleanup, nil
}

func NewRabbitMQChannel(conn *amqp.Connection) (*amqp.Channel, func(), error) {
	logger := slog.With(slog.String("component", componentChannel))

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("failed to open rabbitmq channel", slog.Any("error", err))
		return nil, nil, fmt.Errorf("%w: %v", domain.ErrOpenChannel, err)
	}

	logger.Debug("rabbitmq channel opened")

	cleanup := func() {
		logger.Debug("closing rabbitmq channel")
		if err := ch.Close(); err != nil {
			logger.Error("failed to close rabbitmq channel",
				slog.Any("error", fmt.Errorf("%w: %v", domain.ErrCloseChannel, err)),
			)
		}
	}
	return ch, cleanup, nil
}
func maskRabbitURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "invalid-rabbit-url"
	}

	return u.Redacted()
}
