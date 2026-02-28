package rabbitmq

import (
	"context"
	"encoding/json"

	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeType       = "direct"
	contentTypeJSON    = "application/json"
	componentPublisher = "publisher"
)

type Publisher struct {
	ch     *amqp.Channel
	queue  config.PublisherQueueName
	logger *slog.Logger
}

func NewPublisher(ch *amqp.Channel, queue config.PublisherQueueName) (*Publisher, error) {
	p := &Publisher{
		ch:     ch,
		queue:  queue,
		logger: slog.With(slog.String("component", componentPublisher)),
	}

	if err := p.declareExchange(); err != nil {
		return nil, fmt.Errorf("new publisher: %w", err)
	}

	return p, nil
}

func (p *Publisher) PublishTask(ctx context.Context, taskMessage outPkg.ScrapeTask, routingKey string) error {
	jsonBody, err := json.Marshal(taskMessage)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrMarshalTask, err)
	}

	err = p.ch.PublishWithContext(
		ctx,
		outPkg.RabbitMQExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  contentTypeJSON,
			Body:         jsonBody,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrPublishMessage, err)
	}

	return nil
}

func (p *Publisher) declareExchange() error {
	err := p.ch.ExchangeDeclare(
		outPkg.RabbitMQExchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrDeclareExchange, err)
	}

	return nil
}
