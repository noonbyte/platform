package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func Validation(c fiber.Ctx, err error) error {
	var details []map[string]string
	for _, e := range err.(validator.ValidationErrors) {
		// Customize the error message as needed
		details = append(details, map[string]string{
			"field": e.StructField(),
			"error": fmt.Sprintf("Invalid value for %s", e.Tag()),
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "ERR_VALIDATION",
			"message": "Request validation failed",
			"details": details,
		},
	})
}
