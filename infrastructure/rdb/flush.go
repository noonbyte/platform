package rdb

import (
	"context"
	"fmt"
)

func (r *rdb) FlushKeysByPattern(
	ctx context.Context,
	pattern string,
) error {
	const batchSize = 1000

	var keys []string

	iter := r.client.Scan(ctx, 0, pattern, batchSize).Iterator()

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())

		if len(keys) < batchSize {
			continue
		}

		if err := r.client.Unlink(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("unlink keys: %w", err)
		}

		keys = keys[:0]
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("scan keys: %w", err)
	}

	if len(keys) > 0 {
		if err := r.client.Unlink(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("unlink keys: %w", err)
		}
	}

	return nil
}
