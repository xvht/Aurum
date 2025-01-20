package index

import (
	"github.com/gofiber/fiber/v3"
)

// GET /
func Hello(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Stop poking around!",
		"result":  nil,
	})
}
