package pfcp

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

// InitializeRedis initializes the Redis client
func InitializeRedis(addr string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis")
	}
}

// SaveAssociation saves association information to Redis
func SaveAssociation(nodeID string, data string) error {
	err := redisClient.Set(ctx, "association:"+nodeID, data, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Error saving association: %v", err)
	}
	return err
}

// GetAssociation retrieves association information from Redis
func GetAssociation(nodeID string) (string, error) {
	data, err := redisClient.Get(ctx, "association:"+nodeID).Result()
	if err != nil {
		log.Printf("Error retrieving association: %v", err)
		return "", err
	}
	return data, nil
}

// DeleteAssociation deletes association information from Redis
func DeleteAssociation(nodeID string) error {
	err := redisClient.Del(ctx, "association:"+nodeID).Err()
	if err != nil {
		log.Printf("Error deleting association: %v", err)
	}
	return err
}

