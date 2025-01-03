package src

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
	return &RedisClient{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisClient) SaveSession(sessionID string, session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	return r.client.Set(r.ctx, "session:"+sessionID, data, 0).Err()
}

func (r *RedisClient) GetSession(sessionID string) (*Session, error) {
	data, err := r.client.Get(r.ctx, "session:"+sessionID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	return &session, nil
}

// DeleteSession removes a session from Redis
func (r *RedisClient) DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	if err := r.client.Del(r.ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}

	return nil
}
