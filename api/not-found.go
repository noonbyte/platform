package api

import "github.com/gofiber/fiber/v3"

func NotFound(c fiber.Ctx, code, description string) error {
	return ErrorResponse(c, fiber.StatusNotFound, Error{
		Code:        code,
		Description: description,
	})
}
