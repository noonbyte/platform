package configs

type AppConfiguration struct {
	Name string

	Host string
	Port int

	Environment string

	Domain string
}

func GetAppConfiguration() *AppConfiguration {
	return &AppConfiguration{
		Name: GetEnv("APP_NAME", "backend"),
		Host: GetEnv("APP_HOST", "0.0.0.0"),
		Port: GetEnvAsInt("APP_PORT", 8080),

		Environment: GetEnv("APP_ENVIRONMENT", "production"),

		Domain: GetEnv("APP_DOMAIN", "localhost"),
	}
}
