package bootstrap

import (
	"github.com/noonbyte/platform/configs"
	"github.com/noonbyte/platform/observability/logger"
)

type Config struct {
	AppName string
	Version string

	Logger logger.Config

	Database *configs.DatabaseConfiguration
	Redis    *configs.RedisConfiguration
	NATS     *configs.NATSConfiguration
	S3       *configs.S3Configuration
}
