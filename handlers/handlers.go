package handlers

import (
	"aurum/handlers/v1/misc"

	"github.com/gofiber/fiber/v3"
)

type Route struct {
	Handler     fiber.Handler
	Protected   bool
	Middlewares []fiber.Handler
}

var Handlers = map[string]Route{
	// Default Handlers
	"GET /": {
		Handler:     misc.GetVersion,
		Protected:   false,
		Middlewares: []fiber.Handler{},
	},
	"GET /v1": {
		Handler:     misc.GetVersion,
		Protected:   false,
		Middlewares: []fiber.Handler{},
	},
}
