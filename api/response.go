package api

import "github.com/gofiber/fiber/v3"

func Response(c fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})
}

func MetaResponse(c fiber.Ctx, status int, data any, meta any) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
		"meta": meta,
	})
}

func ErrorResponse(c fiber.Ctx, status int, err Error) error {
	return c.Status(status).JSON(fiber.Map{
		"error": err,
	})
}
