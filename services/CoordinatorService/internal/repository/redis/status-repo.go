package redis

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisStatusRepo struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisTokenRepository(client *redis.Client) domain.StatusRepository {
	return &RedisStatusRepo{client: client, logger: slog.With(slog.String("component", "redisStatusRepository"))}
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
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return id, err
}

func (r *RedisStatusRepo) GetSources(ctx context.Context) ([]outPkg.Source, error) {
	key := outPkg.ScraperRegistryKey
	res, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	sources := make([]outPkg.Source, 0, len(res))
	for _, val := range res {
		var src outPkg.Source
		if err := json.Unmarshal([]byte(val), &src); err != nil {
			continue
		}
		sources = append(sources, src)
	}
	r.logger.Debug("successfully fetched sources from redis", slog.Int("count", len(sources)))
	return sources, nil
}
