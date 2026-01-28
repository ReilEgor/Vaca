package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SubscriberQueueName string
type SubscriberRoutingKey string
type SubscriberExchange string
type TaskSubscriber struct {
	ch         *amqp.Channel
	usecase    domain.ScraperUsecase
	logger     *slog.Logger
	queueName  SubscriberQueueName
	routingKey SubscriberRoutingKey
	exchange   SubscriberExchange
}

func NewTaskSubscriber(
	ch *amqp.Channel,
	uc domain.ScraperUsecase,
	logger *slog.Logger,
	qName SubscriberQueueName,
	rKey SubscriberRoutingKey,
	exch SubscriberExchange,
) *TaskSubscriber {
	return &TaskSubscriber{
		ch:         ch,
		usecase:    uc,
		logger:     logger.With(slog.String("component", "subscriber")),
		queueName:  qName,
		routingKey: rKey,
		exchange:   exch,
	}
}
func (s *TaskSubscriber) Listen(ctx context.Context) error {
	err := s.ch.ExchangeDeclare(
		string(s.exchange),
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = s.ch.QueueDeclare(string(s.queueName), true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = s.ch.QueueBind(string(s.queueName), string(s.routingKey), string(s.exchange), false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := s.ch.Consume(
		string(s.queueName),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	s.logger.Info("started listening for tasks", slog.String("queue", string(s.queueName)))

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("subscriber context cancelled, stopping...")
			return nil
		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("rabbitmq channel closed")
			}

			s.handleDelivery(ctx, d)
		}
	}
}

func (s *TaskSubscriber) handleDelivery(ctx context.Context, d amqp.Delivery) {
	var task outPkg.ScrapeTask

	if err := json.Unmarshal(d.Body, &task); err != nil {
		s.logger.Error("failed to unmarshal task", slog.Any("error", err))
		d.Nack(false, false)
		return
	}

	s.logger.Info("received task", slog.Any("task_id", task.ID))

	err := s.usecase.Execute(ctx, task)
	if err != nil {
		s.logger.Error("failed to execute scrape task",
			slog.Any("task_id", task.ID),
			slog.Any("error", err),
		)
		d.Nack(false, true)
		return
	}

	if err := d.Ack(false); err != nil {
		s.logger.Error("failed to ack message", slog.Any("error", err))
	}
}
