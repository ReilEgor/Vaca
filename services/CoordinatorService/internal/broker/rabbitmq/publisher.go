package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PublisherQueueName string
type Publisher struct {
	ch     *amqp.Channel
	queue  PublisherQueueName
	logger *slog.Logger
}

func NewPublisher(ch *amqp.Channel, queue PublisherQueueName) *Publisher {
	p := &Publisher{
		ch:     ch,
		queue:  queue,
		logger: slog.With(slog.String("component", "publisher")),
	}

	if err := p.DeclareExchange(); err != nil {
		return nil
	}

	return p
}

func (p *Publisher) PublishTask(ctx context.Context, taskMessage outPkg.ScrapeTask, routingKey string) error {
	jsonBody, err := json.Marshal(taskMessage)
	if err != nil {
		p.logger.Error("failed to marshal task message", slog.Any("error", err))
		return fmt.Errorf("marshal error: %w", err)
	}

	return p.ch.PublishWithContext(ctx,
		outPkg.RabbitMQExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonBody,
		})
}

func (p *Publisher) DeclareExchange() error {
	err := p.ch.ExchangeDeclare(
		outPkg.RabbitMQExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}
