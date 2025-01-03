package udr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration int64) error
	Get(ctx context.Context, key string) (string, error)
}

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(addr string) *RedisStorage {
	client := redis.NewClient(&redis.Options{Addr: addr})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	return &RedisStorage{client: client, ctx: context.Background()}
}

func (r *RedisStorage) SaveSubscription(subscription Subscription) error {
	data, err := json.Marshal(subscription)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %w", err)
	}
	return r.client.Set(r.ctx, "subscription:"+subscription.UEID, data, 0).Err()
}

func (r *RedisStorage) GetSubscription(ueID string) (*Subscription, error) {
	data, err := r.client.Get(r.ctx, "subscription:"+ueID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("subscription not found for UE: %s", ueID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve subscription: %w", err)
	}

	var subscription Subscription
	if err := json.Unmarshal([]byte(data), &subscription); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}
	return &subscription, nil
}

