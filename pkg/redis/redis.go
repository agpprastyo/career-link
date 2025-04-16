package redis

import (
	"context"
	"fmt"
	"github.com/agpprastyo/career-link/config"
	"github.com/valyala/fasthttp"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps Redis client functionality
type Client struct {
	client *redis.Client
}

// NewClient creates a new Redis client
func NewClient(cfg config.AppConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Client{client: client}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.client.Close()
}

// Set stores a key-value pair with expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Delete removes keys
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	return result > 0, err
}

// GetClient returns the underlying Redis client
func (c *Client) GetClient() *redis.Client {
	return c.client
}

func (c *Client) Expire(background context.Context, id string, duration time.Duration) {
	c.client.Expire(background, id, duration)
}

func (c *Client) Del(ctx *fasthttp.RequestCtx, id string) interface{} {
	return c.client.Del(ctx, id)
}
