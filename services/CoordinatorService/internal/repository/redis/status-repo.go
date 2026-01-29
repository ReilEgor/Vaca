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

func (r *RedisStatusRepo) Set(ctx context.Context, taskID string, searchKey string, totalSources int, ttl time.Duration) error {
	taskKey := "task:" + taskID

	data := map[string]interface{}{
		"status":    "processing",
		"total":     totalSources,
		"completed": 0,
	}

	if err := r.client.HSet(ctx, taskKey, data).Err(); err != nil {
		return err
	}

	r.client.Expire(ctx, taskKey, ttl)

	hashKey := "hash:" + searchKey
	return r.client.Set(ctx, hashKey, taskID, ttl).Err()
}

func (r *RedisStatusRepo) Get(ctx context.Context, taskID string) map[string]string {
	taskKey := "task:" + taskID
	return r.client.HGetAll(ctx, taskKey).Val()
}

func (r *RedisStatusRepo) GetIDByHash(ctx context.Context, searchKey string) (string, error) {
	hashKey := "hash:" + searchKey
	id, err := r.client.Get(ctx, hashKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	return id, err
}
