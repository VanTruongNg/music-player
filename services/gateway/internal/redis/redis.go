package redis

import (
	"context"
	"fmt"
	"gateway/configs"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *configs.RedisConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       0,
	})
}

func PingRedis(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}
