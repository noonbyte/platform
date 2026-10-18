package configs

type NATSConfiguration struct {
	Host string
	Port int
}

func GetNATSConfiguration() *NATSConfiguration {
	return &NATSConfiguration{
		Host: GetEnv("NATS_HOST", "localhost"),
		Port: GetEnvAsInt("NATS_PORT", 4222),
	}
}
