// internal/redis.go

package internal

import (
	"context"
	"github.com/dedenfarhanhub/blog-service/config"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background()

// RedisClient represents the Redis client
var RedisClient *redis.Client

// InitRedis initializes the Redis connection
func InitRedis() *redis.Client {
	cfg := config.LoadConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort, // Use environment variables
		Password: "",                                  // No password set
		DB:       0,                                   // Use default DB
	})

	// Test the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return client
}
