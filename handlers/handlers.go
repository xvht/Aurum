package handlers

import (
	"vexal/handlers/index"
	"vexal/handlers/misc"

	"github.com/gofiber/fiber/v3"
)

var Handlers = map[string]func(c fiber.Ctx) error{
	// Maintenance Handlers
	"GET /api/version": misc.GetVersion, // Returns API version
	"GET /api/health":  misc.GetHealth,  // Returns API health

	// Default Handler
	"GET /": index.Hello,
}
