package redis

import (
	"context"
	"fmt"

	"github.com/ReilEgor/Vaca/services/StateService/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(host config.Host, port config.RedisPort, password config.Password, db config.DB) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: string(password),
		DB:       int(db),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
