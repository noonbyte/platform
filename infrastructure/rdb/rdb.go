package rdb

import (
	"context"
	"fmt"

	"github.com/noonbyte/platform/configs"
	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Client() *redis.Client
	Close() error

	FlushKeysByPattern(
		ctx context.Context,
		pattern string,
	) error
}

type rdb struct {
	client *redis.Client
}

func New(cfg configs.RedisConfiguration) (Redis, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	if _, err := r.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &rdb{client: r}, nil
}

func (r *rdb) Client() *redis.Client {
	return r.client
}

func (r *rdb) Close() error {
	if r.client == nil {
		return fmt.Errorf("redis: client not initialized")
	}

	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	return nil
}
