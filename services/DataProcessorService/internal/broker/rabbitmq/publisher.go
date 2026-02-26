package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DataPublisher struct {
	ch        *amqp.Channel
	logger    *slog.Logger
	queueName config.PublisherQueueName
}

func NewDataPublisher(
	ch *amqp.Channel,
	qName config.PublisherQueueName,
) *DataPublisher {
	return &DataPublisher{
		ch:        ch,
		logger:    slog.With(slog.String("component", "publisher")),
		queueName: qName,
	}
}

func (p *DataPublisher) Publish(ctx context.Context, message outPkg.ScrapeResult) error {
	jsonBody, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("failed to marshal task message", slog.Any("error", err))
		return fmt.Errorf("marshal error: %w", err)
	}

	return p.ch.PublishWithContext(ctx, "", string(p.queueName), false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         jsonBody,
	})
}
