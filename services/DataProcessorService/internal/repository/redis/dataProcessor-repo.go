package redis

import (
	"context"
	"fmt"

	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisStatusRepo struct {
	client *redis.Client
}

func NewRedisTokenRepository(client *redis.Client) domain.TaskCache {
	return &RedisStatusRepo{client: client}
}

func (r *RedisStatusRepo) IncrementCompleted(ctx context.Context, taskID string) (int64, error) {
	key := "task:" + taskID
	return r.client.HIncrBy(ctx, key, "completed", 1).Result()
}

func (r *RedisStatusRepo) SetStatus(ctx context.Context, taskID string, status string) error {
	taskKey := "task:" + taskID
	err := r.client.HSet(ctx, taskKey, "status", status).Err()
	if err != nil {
		return fmt.Errorf("failed to update status in redis: %w", err)
	}

	return nil
}

func (r *RedisStatusRepo) GetTotal(ctx context.Context, taskID string) (int64, error) {
	taskKey := "task:" + taskID

	val, err := r.client.HGet(ctx, taskKey, "total").Int64()
	if err != nil {
		return 0, err
	}
	return val, nil
}
