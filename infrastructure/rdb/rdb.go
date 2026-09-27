package rdb

import (
	"context"
	"fmt"
	"noonbyte/platform/configs"

	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Client() *redis.Client
	Close() error
}

type redisDB struct {
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

	return &redisDB{client: r}, nil
}

func (r *redisDB) Client() *redis.Client {
	return r.client
}

func (r *redisDB) Close() error {
	if r.client == nil {
		return fmt.Errorf("redis: client not initialized")
	}

	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	return nil
}

func FlushKeysByPattern(ctx context.Context, rdb *redis.Client, pattern string) error {
	// Scan keys by pattern
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		// Delete the key
		err := rdb.Del(ctx, key).Err()
		if err != nil {
			return fmt.Errorf("error deleting key %s: %v", key, err)
		}
		fmt.Println("Deleted key:", key)
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error scanning keys: %v", err)
	}

	return nil
}
