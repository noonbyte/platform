package configs

type JWTConfiguration struct {
	Secret string
}

func GetJWTConfiguration() *JWTConfiguration {
	return &JWTConfiguration{
		Secret: GetEnv("JWT_SECRET", "secret"),
	}
}
