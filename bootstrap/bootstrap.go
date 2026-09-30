package bootstrap

import (
	"context"
	"fmt"

	"github.com/noonbyte/platform/infrastructure/db"
	"github.com/noonbyte/platform/infrastructure/nats"
	"github.com/noonbyte/platform/infrastructure/rdb"
	"github.com/noonbyte/platform/infrastructure/s3"
	"github.com/noonbyte/platform/observability/logger"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, cfg Config) (*Dependencies, error) {
	logger.Init(cfg.Logger)

	log.Info().
		Str("service", cfg.AppName).
		Str("version", cfg.Version).
		Msg("initializing service")

	deps := &Dependencies{}

	if cfg.Database != nil {
		database, err := db.New(*cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("initialize database: %w", err)
		}

		deps.DB = database

		log.Info().
			Msg("database connected")
	}

	if cfg.Redis != nil {
		redis, err := rdb.New(*cfg.Redis)
		if err != nil {
			_ = deps.Close()

			return nil, fmt.Errorf("initialize redis: %w", err)
		}

		deps.Redis = redis

		log.Info().
			Msg("redis connected")
	}

	if cfg.NATS != nil {
		natsClient, err := nats.New(cfg.NATS.Server)
		if err != nil {
			_ = deps.Close()

			return nil, fmt.Errorf("initialize nats: %w", err)
		}

		deps.NATS = natsClient

		log.Info().
			Str("server", cfg.NATS.Server).
			Msg("nats connected")
	}

	if cfg.S3 != nil {
		s3Client, err := s3.New(cfg.S3)
		if err != nil {
			_ = deps.Close()

			return nil, fmt.Errorf("initialize s3: %w", err)
		}

		deps.S3 = s3Client

		log.Info().
			Str("endpoint", cfg.S3.Endpoint).
			Msg("s3 connected")
	}

	log.Info().
		Msg("service dependencies initialized")

	return deps, nil
}
