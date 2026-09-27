package configs

type NATSConfiguration struct {
	Server string
}

func GetNATSConfiguration() *NATSConfiguration {
	return &NATSConfiguration{
		Server: GetEnv("NATS_SERVER", "nats://localhost:4222"),
	}
}
