package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis creates and verifies a Redis client connection.
func ConnectRedis(ctx context.Context, address string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: address})
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}
