package api

import "github.com/gofiber/fiber/v3"

func Forbidden(c fiber.Ctx, code, description string) error {
	return ErrorResponse(c, fiber.StatusForbidden, Error{
		Code:        code,
		Description: description,
	})
}
