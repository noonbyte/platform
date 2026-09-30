package bootstrap

import (
	"github.com/noonbyte/platform/infrastructure/db"
	"github.com/noonbyte/platform/infrastructure/nats"
	"github.com/noonbyte/platform/infrastructure/rdb"
	"github.com/noonbyte/platform/infrastructure/s3"
)

type Dependencies struct {
	DB    db.Database
	Redis rdb.Redis
	NATS  nats.NATS
	S3    s3.S3Client
}

func (d *Dependencies) Close() error {
	var firstErr error

	closeWithError := func(fn func() error) {
		if err := fn(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if d.Redis != nil {
		closeWithError(d.Redis.Close)
	}

	if d.DB != nil {
		closeWithError(d.DB.Close)
	}

	return firstErr
}
