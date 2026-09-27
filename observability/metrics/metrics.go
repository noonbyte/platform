package metrics

import (
	"context"
	"fmt"
	"net/http"
	"noonbyte/platform/configs"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Metrics struct {
	HttpRequestsTotal      *prometheus.CounterVec
	HttpRequestDuration    *prometheus.HistogramVec
	HttpRequestsInProgress prometheus.Gauge
	server                 *http.Server
}

func New() *Metrics {
	return &Metrics{
		HttpRequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		}, []string{"method", "path", "status"}),

		HttpRequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),

		HttpRequestsInProgress: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "Current number of HTTP requests in progress",
		}),
	}
}

func (m *Metrics) StartServer(ctx context.Context, configuration configs.PrometheusConfiguration, log *zap.SugaredLogger) error {
	host := fmt.Sprintf("%s:%d", configuration.ExporterHost, configuration.ExporterPort)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	m.server = &http.Server{
		Addr:    host,
		Handler: mux,
	}

	go func() {
		log.Infof("Metrics server starting on http://%s", host)
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Failed to start metrics server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Info("Shutting down metrics server gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown metrics server gracefully: %w", err)
	}

	log.Info("Metrics server stopped gracefully")
	return nil
}
