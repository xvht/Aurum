package query

import (
	"aurum/packages/crypto/tracker"
	"aurum/types"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	fiberWs "github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	clients    = make(map[*fiberWs.Conn]bool)
	clientsMux sync.RWMutex
)

func WebSocketHandler(c *fiber.Ctx) error {
	if fiberWs.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WebSocketEndpoint(c *fiberWs.Conn) {
	clientsMux.Lock()
	clients[c] = true
	clientsMux.Unlock()

	for {
		var walletQuery types.BalanceQuery
		err := c.ReadJSON(&walletQuery)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close 1001") || strings.Contains(err.Error(), "websocket: close 1006") {
				break
			}

			logrus.Errorf("read error: %v", err)
			break
		}

		queryId := uuid.New().String()

		response := types.BalanceResponse{QueryId: queryId, Error: false, Code: 200}
		walletData, found, err := tracker.GetWalletBalance(strings.ToUpper(walletQuery.Chain), walletQuery.Address)
		if err != nil {
			response.Error = true
			response.Code = 500
		} else if !found {
			response.Error = true
			response.Code = 404
		} else {
			response.Data = walletData
		}

		err = c.WriteJSON(response)
		if err != nil {
			logrus.Errorf("write error: %v", err)
			break
		}
	}

	clientsMux.Lock()
	delete(clients, c)
	clientsMux.Unlock()
}
