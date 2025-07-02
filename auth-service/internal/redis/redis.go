package redis

import (
	"auth-service/configs"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client using RedisConfig.
func NewRedisClient(cfg *configs.RedisConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       0,
	})
}

// PingRedis checks the connection to Redis.
func PingRedis(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}
