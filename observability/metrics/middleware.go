package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

func (m *Metrics) Middleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		m.HttpRequestsInProgress.Inc()
		defer m.HttpRequestsInProgress.Dec()

		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())

		m.HttpRequestsTotal.WithLabelValues(
			c.Method(),
			c.Path(),
			status,
		).Inc()

		m.HttpRequestDuration.WithLabelValues(
			c.Method(),
			c.Path(),
			status,
		).Observe(duration)

		return err
	}
}
