package configs

type RedisConfiguration struct {
	Addr     string
	Password string
	DB       int
}

func GetRedisConfiguration() *RedisConfiguration {
	return &RedisConfiguration{
		Addr:     GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       GetEnvAsInt("REDIS_DB", 0),
	}
}
