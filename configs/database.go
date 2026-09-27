package configs

type DatabaseConfiguration struct {
	Host     string
	Port     uint
	User     string
	Password string
	Name     string
}

func GetDatabaseConfiguration() *DatabaseConfiguration {
	return &DatabaseConfiguration{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     uint(GetEnvAsInt("DB_PORT", 5432)),
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", "password"),
		Name:     GetEnv("DB_NAME", "app"),
	}
}
