package misc

import (
	"aurum/env"

	"github.com/gofiber/fiber/v2"
)

// GET /api/version
func GetVersion(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"code":  200,
		"data": fiber.Map{
			"commit": env.COMMIT_HASH,
			"branch": env.BRANCH,
			"remote": env.REMOTE,
		},
	})
}
