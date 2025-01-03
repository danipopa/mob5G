package pcf

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Storage interface {
	SavePolicy(policy Policy) error
	DeletePolicy(policyId string) error
	GetPolicy(policyId string) (*Policy, error)
	SaveMobilityPolicy(policy AMPolicyResponse) error // Add this method

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

func (r *RedisStorage) SavePolicy(policy Policy) error {
	data, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %w", err)
	}
	return r.client.Set(r.ctx, "policy:"+policy.ID, data, 0).Err()
}

func (r *RedisStorage) DeletePolicy(policyId string) error {
	return r.client.Del(r.ctx, "policy:"+policyId).Err()
}

func (r *RedisStorage) GetPolicy(policyId string) (*Policy, error) {
	data, err := r.client.Get(r.ctx, "policy:"+policyId).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("policy not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	var policy Policy
	if err := json.Unmarshal([]byte(data), &policy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal policy: %w", err)
	}
	return &policy, nil
}

func (r *RedisStorage) SaveMobilityPolicy(policy AMPolicyResponse) error {
	data, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("failed to marshal mobility policy: %w", err)
	}
	return r.client.Set(r.ctx, "mobility_policy:"+policy.PolicyID, data, 0).Err()
}
