package cache

import (
	"context"
	"fmt"

	"github.com/noonbyte/platform/infrastructure/rdb"
	"github.com/noonbyte/platform/specs"
)

func (c *EntityCache[T]) GetList(ctx context.Context, specs specs.Specs) ([]string, error) {
	result, err := c.client.SMembers(ctx, c.ListKey(specs)).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *EntityCache[T]) SetList(
	ctx context.Context,
	specs specs.Specs,
	ids []string,
) error {
	key := c.ListKey(specs)

	if len(ids) == 0 {
		return c.client.Del(ctx, key).Err()
	}

	pipe := c.client.TxPipeline()
	pipe.Del(ctx, key)
	pipe.SAdd(ctx, key, ids)

	ttl := c.ttl

	if len(specs) > 0 {
		ttl = rdb.ScaleTTL(ttl, 10)
	}

	pipe.Expire(ctx, key, c.ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (c *EntityCache[T]) AddToList(
	ctx context.Context,
	id string,
) error {
	key := c.ListKey(specs.NewSpecs())

	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		if err := c.client.SAdd(ctx, key, id).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (c *EntityCache[T]) RemoveFromList(
	ctx context.Context,
	id string,
) error {
	key := c.ListKey(specs.NewSpecs())

	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		if err := c.client.SRem(ctx, key, id).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (c *EntityCache[T]) InvalidateLists(ctx context.Context) error {
	key := fmt.Sprintf("%s:list", c.prefix)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	if err := c.rdb.FlushKeysByPattern(ctx, fmt.Sprintf("%s:*", key)); err != nil {
		return err
	}

	return nil
}
