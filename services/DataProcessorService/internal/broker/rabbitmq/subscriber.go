package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SubscriberQueueName string

type DataSubscriber struct {
	ch        *amqp.Channel
	usecase   domain.DataProcessorUsecase
	logger    *slog.Logger
	queueName SubscriberQueueName
}

func NewTaskSubscriber(
	ch *amqp.Channel,
	uc domain.DataProcessorUsecase,
	logger *slog.Logger,
	qName SubscriberQueueName,
) *DataSubscriber {
	return &DataSubscriber{
		ch:        ch,
		usecase:   uc,
		logger:    logger.With(slog.String("component", "subscriber")),
		queueName: qName,
	}
}

func (s *DataSubscriber) Listen(ctx context.Context) error {
	_, err := s.ch.QueueDeclare(string(s.queueName), true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := s.ch.Consume(string(s.queueName), "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume from queue %s: %w", string(s.queueName), err)
	}

	s.logger.Info("DataProcessor started listening", slog.String("queue", string(s.queueName)))

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("subscriber context cancelled, stopping...")
			return nil
		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("rabbitmq channel closed")
			}

			var result outPkg.ScrapeResult

			if err := json.Unmarshal(d.Body, &result); err != nil {
				s.logger.Error("failed to unmarshal result", slog.Any("error", err))
				d.Nack(false, false)
				continue
			}

			err := s.usecase.Process(ctx, result)
			if err != nil {
				s.logger.Error("failed to process vacancies",
					slog.Any("task_id", result.TaskID),
					slog.Any("error", err))

				d.Nack(false, true)
				continue
			}

			if err := d.Ack(false); err != nil {
				s.logger.Error("failed to ack message", slog.Any("error", err))
			}
		}
	}
}
