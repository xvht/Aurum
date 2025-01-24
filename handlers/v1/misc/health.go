package misc

import (
	"github.com/gofiber/fiber/v2"
)

// GET /v1/health
func GetHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"code":  200,
		"data":  "API is healthy",
	})
}
