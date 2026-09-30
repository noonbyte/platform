package validators

import (
	"github.com/gofiber/fiber/v3"
	"github.com/noonbyte/platform/api"
)

func isID(c fiber.Ctx, id string) error {
	if len(id) != 32 {
		return api.BadRequest(c, "ERR_VALIDATION", "Id is required")
	}

	return nil
}
