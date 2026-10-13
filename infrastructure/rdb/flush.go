package rdb

import (
	"context"
	"fmt"
	"path"
)

func (r *rdb) FlushKeysByPattern(
	ctx context.Context,
	pattern string,
	excludePatterns ...string,
) error {
	const batchSize = 1000

	var keys []string

	iter := r.client.Scan(ctx, 0, pattern, batchSize).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()

		excluded := false
		for _, excludePattern := range excludePatterns {
			matched, err := path.Match(excludePattern, key)
			if err != nil {
				return fmt.Errorf("match exclude pattern %q: %w", excludePattern, err)
			}
			if matched {
				excluded = true
				break
			}
		}

		if excluded {
			continue
		}

		keys = append(keys, key)

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
