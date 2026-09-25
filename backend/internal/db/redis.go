package db

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context, addr string) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{Addr: addr})
	if err := r.Ping(ctx).Err(); err != nil {
		_ = r.Close()
		return nil, err
	}
	return r, nil
}
