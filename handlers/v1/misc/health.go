package misc

import (
	"github.com/gofiber/fiber/v2"
)

// GET /api/health
func GetHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "API health",
		"result":  "healthy",
	})
}
