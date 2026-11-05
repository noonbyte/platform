package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

func (m *Metrics) Middleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		m.IncRequestsInProgress()
		defer m.DecRequestsInProgress()

		err := c.Next()

		status := c.Response().StatusCode()
		if err != nil {
			status = fiber.StatusInternalServerError
		}

		path := c.Route().Path
		if path == "" {
			path = "unmatched"
		}

		m.RecordRequest(
			c.Method(),
			path,
			status,
			time.Since(start),
			int64(len(c.Request().Body())),
			int64(len(c.Response().Body())),
		)

		return err
	}
}

func (m *Metrics) StatusCode(status int) string {
	return strconv.Itoa(status)
}
