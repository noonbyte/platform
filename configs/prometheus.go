package configs

type PrometheusConfiguration struct {
	ExporterHost string
	ExporterPort int
}

func GetPrometheusConfiguration() *PrometheusConfiguration {
	return &PrometheusConfiguration{
		ExporterHost: GetEnv("PROMETHEUS_EXPORTER_HOST", "0.0.0.0"),
		ExporterPort: GetEnvAsInt("PROMETHEUS_EXPORTER_PORT", 9100),
	}
}
