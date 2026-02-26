package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/NotificationService/internal/config"
	"github.com/ReilEgor/Vaca/services/NotificationService/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationSubscriber struct {
	ch        *amqp.Channel
	logger    *slog.Logger
	usecase   domain.NotificationUsecase
	queueName config.SubscriberQueueName
}

func NewNotificationSubscriber(
	ch *amqp.Channel,
	uc domain.NotificationUsecase,
	qName config.SubscriberQueueName,
) *NotificationSubscriber {
	return &NotificationSubscriber{
		ch:        ch,
		logger:    slog.With(slog.String("component", "notificationSubscriber")),
		usecase:   uc,
		queueName: qName,
	}
}

func (s *NotificationSubscriber) Listen(ctx context.Context) error {
	_, err := s.ch.QueueDeclare(string(s.queueName), true, false, false, false, nil)
	if err != nil {
		s.logger.Error(domain.FailedToDeclareQueue.Error(),
			slog.String("queueName", string(s.queueName)),
			slog.Any("error", err),
		)
		return err
	}
	msgs, err := s.ch.Consume(string(s.queueName), "", false, false, false, false, nil)
	if err != nil {
		s.logger.Error(domain.FailedToConsumeFromQueue.Error(),
			slog.String("queueName", string(s.queueName)),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Debug("started listening", slog.String("queue", string(s.queueName)))

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("stopping subscriber due to context cancellation")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				s.logger.Warn("message channel closed, stopping subscriber")
				return nil
			}

			var message outPkg.ScrapeResult
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				s.logger.Error("failed to unmarshal message",
					slog.Any("error", err),
				)

				err := msg.Nack(false, false)
				if err != nil {
					s.logger.Error(domain.FailedToNeckMessages.Error(), slog.Any("error", err))
					return fmt.Errorf("%w:%v", domain.FailedToNeckMessages, err)
				}

				continue
			}
			err := s.usecase.SendNotification(ctx, message)
			if err != nil {
				return err
			}
			if err := msg.Ack(false); err != nil {
				s.logger.Error("failed to ack message", slog.Any("error", err))
			}
		}
	}
}
