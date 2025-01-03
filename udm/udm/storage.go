package udm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Storage interface {
	GetSubscriptionData(ueID string) (*SubscriptionData, error)
	GetAuthVector(ueID string) (*AuthVector, error)
}

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStorage initializes Redis storage
func NewRedisStorage(addr string) *RedisStorage {
	client := redis.NewClient(&redis.Options{Addr: addr})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	return &RedisStorage{
		client: client,
		ctx:    context.Background(),
	}
}

// GetSubscriptionData retrieves subscription data for a UE from Redis
func (r *RedisStorage) GetSubscriptionData(ueID string) (*SubscriptionData, error) {
	key := fmt.Sprintf("subscription:%s", ueID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("subscription data not found for UE: %s", ueID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve subscription data: %w", err)
	}

	var subscription SubscriptionData
	if err := json.Unmarshal([]byte(data), &subscription); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription data: %w", err)
	}
	return &subscription, nil
}

// GetAuthVector retrieves the authentication vector for a UE from Redis
func (r *RedisStorage) GetAuthVector(ueID string) (*AuthVector, error) {
	key := fmt.Sprintf("auth:%s", ueID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("authentication vector not found for UE: %s", ueID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve authentication vector: %w", err)
	}

	var vector AuthVector
	if err := json.Unmarshal([]byte(data), &vector); err != nil {
		return nil, fmt.Errorf("failed to unmarshal authentication vector: %w", err)
	}
	return &vector, nil
}

