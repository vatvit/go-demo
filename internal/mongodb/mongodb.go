package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func New(uri, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Client{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, nil)
}

func (c *Client) Database() *mongo.Database {
	return c.db
}
