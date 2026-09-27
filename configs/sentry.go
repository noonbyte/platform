package configs

type SentryConfiguration struct {
	DSN string
}

func GetSentryConfiguration() *SentryConfiguration {
	return &SentryConfiguration{
		DSN: GetEnv("SENTRY_DSN", ""),
	}
}
