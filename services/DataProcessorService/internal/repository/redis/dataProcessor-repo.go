package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisStatusRepo struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisTokenRepository(client *redis.Client) domain.TaskCache {
	return &RedisStatusRepo{client: client, logger: slog.With(slog.String("component", "redisStatusRepository"))}
}

func (r *RedisStatusRepo) IncrementCompleted(ctx context.Context, taskID string) (int64, error) {
	key := "task:" + taskID
	val, err := r.client.HIncrBy(ctx, key, "completed", 1).Result()
	if err != nil {
		r.logger.Debug(domain.FailedToIncrementCompleted.Error(), slog.String("task_id", taskID), slog.Any("error", err))
		return 0, fmt.Errorf("%w:%v", domain.FailedToIncrementCompleted, err)
	}
	return val, nil
}

func (r *RedisStatusRepo) SetStatus(ctx context.Context, taskID string, status string) error {
	taskKey := "task:" + taskID
	err := r.client.HSet(ctx, taskKey, "status", status).Err()
	if err != nil {
		r.logger.Info(domain.FailedToUpdateStatus.Error(), slog.String("task_id", taskID), slog.Any("error", err))
		return fmt.Errorf("%w:%v", domain.FailedToUpdateStatus, err)
	}
	return nil
}

func (r *RedisStatusRepo) GetTotal(ctx context.Context, taskID string) (int64, error) {
	taskKey := "task:" + taskID

	val, err := r.client.HGet(ctx, taskKey, "total").Int64()
	if err != nil {
		r.logger.Info(domain.FailedToGetTotal.Error(), slog.String("task_id", taskID), slog.Any("error", err))
		return 0, fmt.Errorf("%w:%v", domain.FailedToGetTotal, err)
	}
	return val, nil
}
