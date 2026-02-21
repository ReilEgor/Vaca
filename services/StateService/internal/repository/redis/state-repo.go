package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/StateService/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisStateRepo struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisTokenRepository(client *redis.Client) domain.StateRepository {
	return &RedisStateRepo{client: client, logger: slog.With(slog.String("component", "redisStateRepository"))}
}

func (r *RedisStateRepo) Set(ctx context.Context, taskID string, searchKey string, totalSources int, ttl time.Duration) error {
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

func (r *RedisStateRepo) Get(ctx context.Context, taskID string) map[string]string {
	taskKey := "task:" + taskID
	return r.client.HGetAll(ctx, taskKey).Val()
}

func (r *RedisStateRepo) GetIDByHash(ctx context.Context, searchKey string) (string, error) {
	hashKey := "hash:" + searchKey
	id, err := r.client.Get(ctx, hashKey).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return id, err
}

func (r *RedisStateRepo) GetSources(ctx context.Context) ([]outPkg.Source, error) {
	key := outPkg.ScraperRegistryKey
	res, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	sources := make([]outPkg.Source, 0, len(res))
	for _, val := range res {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var src outPkg.Source
		if err := json.Unmarshal([]byte(val), &src); err != nil {
			continue
		}
		sources = append(sources, src)
	}
	r.logger.Debug("successfully fetched sources from redis", slog.Int("count", len(sources)))
	return sources, nil
}

func (r *RedisStateRepo) SetStatus(ctx context.Context, taskID string, status string) error {
	taskKey := "task:" + taskID
	err := r.client.HSet(ctx, taskKey, "status", status).Err()
	if err != nil {
		r.logger.Info(domain.FailedToUpdateStatus.Error(), slog.String("task_id", taskID), slog.Any("error", err))
		return fmt.Errorf("%w:%v", domain.FailedToUpdateStatus, err)
	}
	return nil
}

func (r *RedisStateRepo) GetTotal(ctx context.Context, taskID string) (int64, error) {
	taskKey := "task:" + taskID

	val, err := r.client.HGet(ctx, taskKey, "total").Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	} else if err != nil {
		r.logger.Info(domain.FailedToGetTotal.Error(), slog.String("task_id", taskID), slog.Any("error", err))
		return 0, fmt.Errorf("%w:%v", domain.FailedToGetTotal, err)
	}
	return val, nil
}

func (r *RedisStateRepo) IncrementCompleted(ctx context.Context, taskID string) (int64, error) {
	key := "task:" + taskID
	val, err := r.client.HIncrBy(ctx, key, "completed", 1).Result()
	if err != nil {
		r.logger.Debug(domain.FailedToIncrementCompleted.Error(), slog.String("task_id", taskID), slog.Any("error", err))
		return 0, fmt.Errorf("%w:%v", domain.FailedToIncrementCompleted, err)
	}
	return val, nil
}

func (r *RedisStateRepo) Register(ctx context.Context, source outPkg.Source) error {
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
