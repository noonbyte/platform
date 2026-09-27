package api

import "github.com/gofiber/fiber/v3"

func Conflict(c fiber.Ctx, code, description string) error {
	return ErrorResponse(c, fiber.StatusConflict, Error{
		Code:        code,
		Description: description,
	})
}
