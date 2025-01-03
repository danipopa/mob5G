package ausf

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Storage interface {
	SaveAuthSession(ueID string, authData AuthResponse) error
	GetAuthSession(ueID string) (*AuthResponse, error)
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

// SaveAuthSession saves authentication session data for a UE
func (r *RedisStorage) SaveAuthSession(ueID string, authData AuthResponse) error {
	data, err := json.Marshal(authData)
	if err != nil {
		return fmt.Errorf("failed to marshal auth session data: %w", err)
	}
	return r.client.Set(r.ctx, "auth:"+ueID, data, 0).Err()
}

// GetAuthSession retrieves authentication session data for a UE
func (r *RedisStorage) GetAuthSession(ueID string) (*AuthResponse, error) {
	data, err := r.client.Get(r.ctx, "auth:"+ueID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("authentication session not found for UE: %s", ueID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve auth session data: %w", err)
	}

	var authData AuthResponse
	if err := json.Unmarshal([]byte(data), &authData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal auth session data: %w", err)
	}
	return &authData, nil
}

