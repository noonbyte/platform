package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func Logging() fiber.Handler {
	logger := log.Logger

	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		status := c.Response().StatusCode()

		event := logger.Info()

		switch {
		case err != nil || status >= 500:
			event = logger.Error()

			if err != nil {
				event = event.Err(err)
			}

		case status >= 400:
			event = logger.Warn()
		}

		event.
			Str("request_id", c.GetRespHeader("X-Request-ID")).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Int("status", status).
			Int64("latency_ms", time.Since(start).Milliseconds()).
			Msg("http request")

		return err
	}
}
