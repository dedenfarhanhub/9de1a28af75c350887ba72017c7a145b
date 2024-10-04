package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisService struct
type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService initializes redis service
func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		client: client,
		ctx:    context.Background(),
	}
}

// SetEntity sets an entity in Redis with serialization
func (r *RedisService) SetEntity(entityType string, id string, entity interface{}, expiration time.Duration) error {
	entityJSON, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, entityType+":"+id, entityJSON, expiration).Err()
}

// GetEntity retrieves any entity from Redis and deserializes it into the specified type
func (r *RedisService) GetEntity(entityType string, id string, entity interface{}) error {
	entityJSON, err := r.client.Get(r.ctx, entityType+":"+id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil // Entity not found
		}
		return err
	}
	return json.Unmarshal([]byte(entityJSON), entity)
}

// DeleteEntity removes an entity from Redis using its type and ID
func (r *RedisService) DeleteEntity(entityType string, id string) error {
	ctx := context.Background()  // Create a context for the Redis operations
	key := entityType + ":" + id // Construct the key based on entity type and ID

	// Use the Redis DEL command to remove the key
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err // Return the error if deletion fails
	}

	return nil // Return nil if deletion was successful
}
