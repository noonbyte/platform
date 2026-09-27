package api

import "github.com/gofiber/fiber/v3"

func NoContent(c fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}
