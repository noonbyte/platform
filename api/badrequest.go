package api

import "github.com/gofiber/fiber/v3"

func BadRequest(c fiber.Ctx, code, description string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, Error{
		Code:        code,
		Description: description,
	})
}
