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

func (r *RedisStatusRepo) Set(ctx context.Context, taskID string, searchKey, status string, ttl time.Duration) error {
	taskKey := "task:" + taskID
	if err := r.client.Set(ctx, taskKey, status, ttl).Err(); err != nil {
		return err
	}
	hashKey := "hash:" + searchKey
	return r.client.Set(ctx, hashKey, taskID, ttl).Err()
}

func (r *RedisStatusRepo) Get(ctx context.Context, taskID string) (string, error) {
	return r.client.Get(ctx, taskID).Result()
}
func (r *RedisStatusRepo) GetIDByHash(ctx context.Context, searchKey string) (string, error) {
	hashKey := "hash:" + searchKey
	id, err := r.client.Get(ctx, hashKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	return id, err
}
