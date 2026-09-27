package api

import "github.com/gofiber/fiber/v3"

func Unauthorized(c fiber.Ctx, reason string) error {
	if reason != "" {
		return ErrorResponse(c, fiber.StatusUnauthorized, Error{
			Code:        "ERR_UNAUTHORIZED",
			Description: reason,
		})
	}

	return ErrorResponse(c, fiber.StatusUnauthorized, Error{
		Code:        "ERR_UNAUTHORIZED",
		Description: "Authorization failed",
	})
}
