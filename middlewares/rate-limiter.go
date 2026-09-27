package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/noonbyte/platform/api"
)

func RateLimiter(max int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: duration,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return api.ErrorResponse(c, fiber.StatusTooManyRequests, api.Error{
				Code:        "ERR_RATE_LIMITED",
				Description: "Too many requests. Please try again later.",
			})
		},
	})
}
