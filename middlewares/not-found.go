package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/noonbyte/platform/api"
)

func NotFoundMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		return api.NotFound(c, "ERR_INVALID_ROUTE", "Requested route was not found")
	}
}
