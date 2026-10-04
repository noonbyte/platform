package api

import "github.com/gofiber/fiber/v3"

func Created(c fiber.Ctx, id string) error {
	return c.Status(fiber.StatusCreated).JSON(map[string]string{"id": id})

}
