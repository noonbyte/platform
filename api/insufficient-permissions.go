package api

import (
	"github.com/gofiber/fiber/v3"
)

func InsufficientPermissions(c fiber.Ctx) error {
	return ErrorResponse(c, fiber.StatusForbidden, Error{
		Code:        "ERR_INSUFFICIENT_PERMISSIONS",
		Description: "you don't have sufficient permission",
	})
}
