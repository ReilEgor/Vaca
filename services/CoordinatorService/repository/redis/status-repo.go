package redis

import (
	"context"
	"time"

	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisStatusRepo struct {
	client *redis.Client
}

func NewRedisTokenRepository(client *redis.Client) domain.StatusRepository {
	return &RedisStatusRepo{client: client}
}

func (r *RedisStatusRepo) Set(ctx context.Context, taskID string, status string, ttl time.Duration) error {
	return r.client.Set(ctx, taskID, status, ttl).Err()
}

func (r *RedisStatusRepo) Get(ctx context.Context, taskID string) (string, error) {
	return r.client.Get(ctx, taskID).Result()
}
