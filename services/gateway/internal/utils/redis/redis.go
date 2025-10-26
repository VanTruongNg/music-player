package redisutil

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct {
	client *redis.Client
}

// NewRedisUtil creates a new RedisUtil instance with the provided redis.Client.
func NewRedisUtil(client *redis.Client) *RedisUtil {
	return &RedisUtil{client: client}
}

// Set stores a string value with a TTL in Redis.
func (r *RedisUtil) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a string value from Redis by key.
func (r *RedisUtil) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete removes a key from Redis.
func (r *RedisUtil) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// SetJSON marshals a value to JSON and stores it in Redis with a TTL.
func (r *RedisUtil) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

// GetJSON retrieves a JSON value from Redis and unmarshals it into dest.
func (r *RedisUtil) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
