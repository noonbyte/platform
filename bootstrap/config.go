package bootstrap

import (
	"github.com/noonbyte/platform/configs"
)

type Config struct {
	AppName string
	Version string

	Database *configs.DatabaseConfiguration
	Redis    *configs.RedisConfiguration
	NATS     *configs.NATSConfiguration
	S3       *configs.S3Configuration
}
