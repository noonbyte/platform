package api

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func InternalServerError(c fiber.Ctx, err error) error {
	fmt.Println(err.Error())

	return ErrorResponse(c, fiber.StatusInternalServerError, Error{
		Code:        "ERR_INTERNAL_ERROR",
		Description: "Something went wrong! Please try again later.",
	})
}
