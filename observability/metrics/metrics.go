package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/noonbyte/platform/configs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Metrics struct {
	// HTTP metrics
	HttpRequestsTotal       *prometheus.CounterVec
	HttpRequestsErrorsTotal *prometheus.CounterVec
	HttpRequestDuration     *prometheus.HistogramVec
	HttpRequestsInProgress  prometheus.Gauge
	HttpRequestSizeBytes    prometheus.Histogram
	HttpResponseSizeBytes   prometheus.Histogram

	// Application metrics
	ActiveConnections prometheus.Gauge
	ApplicationInfo   *prometheus.GaugeVec
	ApplicationUp     prometheus.Gauge

	// Exporter metrics
	ExporterRequestsTotal prometheus.Counter
	ExporterDuration      prometheus.Histogram

	registry    *prometheus.Registry
	server      *http.Server
	stopRuntime chan struct{}
	runtimeOnce sync.Once
}

func New() *Metrics {
	registry := prometheus.NewRegistry()

	m := &Metrics{
		registry:    registry,
		stopRuntime: make(chan struct{}),

		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),

		HttpRequestsErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_errors_total",
				Help: "Total number of HTTP requests resulting in 4xx or 5xx responses.",
			},
			[]string{"method", "path", "status"},
		),

		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "HTTP request duration in seconds.",
				Buckets: []float64{
					0.005,
					0.01,
					0.025,
					0.05,
					0.1,
					0.25,
					0.5,
					1,
					2.5,
					5,
					10,
				},
			},
			[]string{"method", "path", "status"},
		),

		HttpRequestsInProgress: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_progress",
				Help: "Current number of HTTP requests in progress.",
			},
		),

		HttpRequestSizeBytes: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "HTTP request body size in bytes.",
				Buckets: prometheus.ExponentialBuckets(128, 4, 8),
			},
		),

		HttpResponseSizeBytes: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response body size in bytes.",
				Buckets: prometheus.ExponentialBuckets(128, 4, 8),
			},
		),

		ActiveConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "application_active_connections",
				Help: "Current number of active application connections.",
			},
		),

		ApplicationInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "application_info",
				Help: "Application build and version information.",
			},
			[]string{"version", "environment"},
		),

		ApplicationUp: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "application_up",
				Help: "Application availability status.",
			},
		),

		ExporterRequestsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "metrics_exporter_requests_total",
				Help: "Total number of requests to the metrics exporter.",
			},
		),

		ExporterDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "metrics_exporter_request_duration_seconds",
				Help:    "Metrics exporter request duration.",
				Buckets: prometheus.DefBuckets,
			},
		),
	}

	registry.MustRegister(
		m.HttpRequestsTotal,
		m.HttpRequestsErrorsTotal,
		m.HttpRequestDuration,
		m.HttpRequestsInProgress,
		m.HttpRequestSizeBytes,
		m.HttpResponseSizeBytes,
		m.ActiveConnections,
		m.ApplicationInfo,
		m.ApplicationUp,
		m.ExporterRequestsTotal,
		m.ExporterDuration,

		// Go runtime metrics.
		collectors.NewGoCollector(),

		// Process metrics.
		collectors.NewProcessCollector(
			collectors.ProcessCollectorOpts{},
		),
	)

	return m
}

func (m *Metrics) SetApplicationInfo(version, environment string) {
	m.ApplicationInfo.Reset()
	m.ApplicationInfo.WithLabelValues(version, environment).Set(1)
}

func (m *Metrics) SetApplicationUp(up bool) {
	if up {
		m.ApplicationUp.Set(1)
		return
	}

	m.ApplicationUp.Set(0)
}

func (m *Metrics) RecordRequest(
	method string,
	path string,
	status int,
	duration time.Duration,
	requestSize int64,
	responseSize int64,
) {
	statusCode := fmt.Sprintf("%d", status)

	m.HttpRequestsTotal.WithLabelValues(
		method,
		path,
		statusCode,
	).Inc()

	if status >= 400 {
		m.HttpRequestsErrorsTotal.WithLabelValues(
			method,
			path,
			statusCode,
		).Inc()
	}

	m.HttpRequestDuration.WithLabelValues(
		method,
		path,
		statusCode,
	).Observe(duration.Seconds())

	m.HttpRequestSizeBytes.Observe(float64(requestSize))
	m.HttpResponseSizeBytes.Observe(float64(responseSize))
}

func (m *Metrics) IncRequestsInProgress() {
	m.HttpRequestsInProgress.Inc()
}

func (m *Metrics) DecRequestsInProgress() {
	m.HttpRequestsInProgress.Dec()
}

func (m *Metrics) IncActiveConnections() {
	m.ActiveConnections.Inc()
}

func (m *Metrics) DecActiveConnections() {
	m.ActiveConnections.Dec()
}

func (m *Metrics) StartRuntimeMetrics(ctx context.Context) {
	m.runtimeOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-m.stopRuntime:
					return
				case <-ticker.C:
					var stats runtime.MemStats
					runtime.ReadMemStats(&stats)

					_ = stats
				}
			}
		}()
	})
}

func (m *Metrics) StartServer(
	ctx context.Context,
	configuration configs.PrometheusConfiguration,
) error {
	host := fmt.Sprintf(
		"%s:%d",
		configuration.ExporterHost,
		configuration.ExporterPort,
	)

	mux := http.NewServeMux()

	handler := promhttp.HandlerFor(
		m.registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	)

	mux.Handle("/metrics", handler)

	m.server = &http.Server{
		Addr:              host,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Str("host", host).Msg("Metrics server starting")

		if err := m.server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-serverErrors:
		return fmt.Errorf("metrics server failed: %w", err)
	}

	log.Info().Msg("Shutting down metrics server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	close(m.stopRuntime)

	if err := m.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf(
			"failed to shutdown metrics server gracefully: %w",
			err,
		)
	}

	log.Info().Msg("Metrics server stopped gracefully")

	return nil
}
