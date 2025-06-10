package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/travboz/backend-projects/todo-list-api/internal/env"
)

// RedisConfig holds the configuration for Redis connection
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Timeout  time.Duration
}

// DefaultRedisConfig returns a default Redis configuration
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host: "localhost",
		Port: env.GetString("CACHE_ACCESS_PORT", "6379"),
	}
}

// ConnectRedis configures, connects to Redis, and tests the connection
func NewRedisClient(config *RedisConfig) (*redis.Client, error) {
	// Create Redis client with configuration
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	// Create context with timeout for testing connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the connection with PING command
	_, err := client.Ping(ctx).Result()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
