package configs

type RedisConfiguration struct {
	Host     string
	Port     int
	Password string
	Database int
}

func GetRedisConfiguration() *RedisConfiguration {
	return &RedisConfiguration{
		Host:     GetEnv("REDIS_HOST", "localhost"),
		Port:     GetEnvAsInt("REDIS_PORT", 6379),
		Password: GetEnv("REDIS_PASSWORD", ""),
		Database: GetEnvAsInt("REDIS_DB", 0),
	}
}
