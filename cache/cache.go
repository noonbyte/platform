package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/noonbyte/platform/infrastructure/rdb"
	"github.com/redis/go-redis/v9"
)

type EntityCache[T any] struct {
	rdb    rdb.Redis
	client redis.Client

	prefix string
	ttl    time.Duration
	getID  func(*T) string
}

func NewEntity[T any](
	rdb rdb.Redis,
	prefix string,
	ttl time.Duration,
	getID func(*T) string,
) *EntityCache[T] {
	return &EntityCache[T]{
		rdb:    rdb,
		client: *rdb.Client(),

		prefix: prefix,
		ttl:    ttl,
		getID:  getID,
	}
}

func (c *EntityCache[T]) Get(
	ctx context.Context,
	id string,
) (*T, error) {
	value, err := c.client.Get(ctx, c.Key(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrMiss
		}

		return nil, err
	}

	var entity T

	if err := json.Unmarshal([]byte(value), &entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (c *EntityCache[T]) GetMany(
	ctx context.Context,
	ids []string,
	loader func(ctx context.Context, ids []string) ([]*T, error),
) ([]*T, error) {
	if len(ids) == 0 {
		return []*T{}, nil
	}

	keys := make([]string, len(ids))

	for i, id := range ids {
		keys[i] = c.Key(id)
	}

	values, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make([]*T, 0, len(ids))
	missing := make([]string, 0)

	for i, value := range values {
		if value == nil {
			missing = append(missing, ids[i])
			continue
		}

		raw, ok := value.(string)
		if !ok {
			missing = append(missing, ids[i])
			continue
		}

		var entity T

		if err := json.Unmarshal([]byte(raw), &entity); err != nil {
			missing = append(missing, ids[i])
			continue
		}

		result = append(result, &entity)
	}

	if len(missing) == 0 {
		return result, nil
	}

	entities, err := loader(ctx, missing)
	if err != nil {
		return nil, err
	}

	for _, entity := range entities {
		if entity == nil {
			continue
		}

		result = append(result, entity)

		raw, err := json.Marshal(entity)
		if err != nil {
			return nil, err
		}

		if err := c.client.Set(
			ctx,
			c.Key(c.ID(entity)),
			raw,
			c.ttl,
		).Err(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *EntityCache[T]) Set(
	ctx context.Context,
	entity *T,
) error {
	if entity == nil {
		return nil
	}

	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	return c.client.Set(
		ctx,
		c.Key(c.ID(entity)),
		data,
		c.ttl,
	).Err()
}

func (c *EntityCache[T]) Delete(
	ctx context.Context,
	id string,
) error {
	return c.client.Del(ctx, c.Key(id)).Err()
}
