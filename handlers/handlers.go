package handlers

import (
	"aurum/handlers/v1/misc"
	"aurum/handlers/v1/ws/prices"
	"aurum/handlers/v1/ws/query"

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
		Handler: misc.GetVersion,
	},
	"GET /v1": {
		Handler: misc.GetVersion,
	},
	"GET /v1/health": {
		Handler: misc.GetHealth,
	},

	// WebSocket Handlers
	"SOCKET /v1/ws/prices": {
		Handler:  prices.WebSocketHandler,
		Endpoint: prices.WebSocketEndpoint,
	},
	"SOCKET /v1/ws/query": {
		Handler:  query.WebSocketHandler,
		Endpoint: query.WebSocketEndpoint,
	},
}
