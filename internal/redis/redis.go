package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func New(addr string) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Client() *redis.Client {
	return c.client
}
