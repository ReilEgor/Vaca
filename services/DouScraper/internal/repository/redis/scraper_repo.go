package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/redis/go-redis/v9"
)

type RedisScraperRepo struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisScraperRepo(client *redis.Client) *RedisScraperRepo {
	return &RedisScraperRepo{
		client: client,
		logger: slog.With(slog.String("component", "redisScraperRepo")),
	}
}

func (r *RedisScraperRepo) Register(ctx context.Context, source outPkg.Source, ttl time.Duration) error {
	key := outPkg.ScraperRegistryKey

	data, err := json.Marshal(source)
	if err != nil {
		r.logger.Error("failed to marshal source", slog.Any("error", err))
		return fmt.Errorf("failed to marshal source: %w", err)
	}

	err = r.client.HSet(ctx, key, source.Name, data).Err()
	if err != nil {
		r.logger.Error("failed to register in redis", slog.Any("error", err))
		return fmt.Errorf("failed to register in redis: %w", err)
	}
	r.logger.Debug("successfully registered in redis")
	return nil
}
