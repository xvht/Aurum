package handlers

import (
	"aurum/handlers/v1/crypto"
	"aurum/handlers/v1/misc"

	"github.com/gofiber/fiber/v2"
	fiberWs "github.com/gofiber/websocket/v2"
)

type Route struct {
	Handler     fiber.Handler
	Endpoint    func(*fiberWs.Conn)
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

	// WebSocket Handlers
	"GET /v1/ws/prices": {
		Handler:     crypto.WebSocketHandler,
		Protected:   false,
		Middlewares: []fiber.Handler{},
	},
	"SOCKET /v1/ws/prices": {
		Handler:     crypto.WebSocketHandler,
		Endpoint:    crypto.WebSocketEndpoint,
		Protected:   false,
		Middlewares: []fiber.Handler{},
	},
}
