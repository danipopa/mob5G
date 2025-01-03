package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{Addr: addr})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	return &RedisClient{client: client, ctx: context.Background()}
}

func (r *RedisClient) Save(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return r.client.Set(r.ctx, key, data, 0).Err()
}

func (r *RedisClient) Get(key string, value interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get data: %w", err)
	}
	return json.Unmarshal([]byte(data), value)
}

