package go_redis

import (
	"context"
	"fmt"
	"github.com/MentalMentos/techFin/internal/clients/redis"
	"github.com/goccy/go-json"
	goRedis "github.com/redis/go-redis/v9"
	"time"
)

// GoRedisClient is a wrapper around the go-redis client, providing methods
// for interacting with a Redis data store.
type GoRedisClient struct {
	Client *goRedis.Client // Underlying go-redis client instance
}

// NewGoRedisClient initializes a new GoRedisClient with the provided Redis configuration.
// It creates a new Redis client using the configuration's address.
func NewGoRedisClient(config redis.IRedisConfig) (*GoRedisClient, error) {
	if config == nil {
		return nil, fmt.Errorf("Redis config is nil")
	}

	address := config.Address()
	if address == "" {
		return nil, fmt.Errorf("Redis address is empty")
	}

	// Creating Redis client
	redisClient := &GoRedisClient{
		Client: goRedis.NewClient(&goRedis.Options{
			Addr:     address,
			Password: "", // No password for the default config
			DB:       0,  // Default DB
		}),
	}

	// Checking if Redis is reachable
	_, err := redisClient.Client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s: %v", address, err)
	}

	return redisClient, nil
}

// Set stores the string value in Redis with an expiration time.
func (g *GoRedisClient) Set(ctx context.Context, key string, value string, duration time.Duration) error {
	return g.Client.Set(ctx, key, value, duration).Err()
}

// Get retrieves the string value from Redis.
func (g *GoRedisClient) Get(ctx context.Context, key string) (string, error) {
	return g.Client.Get(ctx, key).Result()
}

// SetObject stores a serialized object in Redis with an expiration time.
func (g *GoRedisClient) SetObject(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	noteBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %v", err)
	}
	return g.Client.Set(ctx, key, string(noteBytes), duration).Err()
}

// GetObject retrieves and unmarshals an object from Redis.
func (g *GoRedisClient) GetObject(ctx context.Context, key string, value any) error {
	val, err := g.Client.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get object from Redis: %v", err)
	}

	if err = json.Unmarshal([]byte(val), &value); err != nil {
		return fmt.Errorf("failed to unmarshal object: %v", err)
	}
	return nil
}

// Delete removes a key from Redis.
func (g *GoRedisClient) Delete(ctx context.Context, key string) error {
	return g.Client.Del(ctx, key).Err()
}
