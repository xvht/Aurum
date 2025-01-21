package query

import (
	"aurum/packages/crypto/tracker"
	"aurum/types"
	"strings"
	"sync"
	"time"

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

const pingInterval = 30 * time.Second

func WebSocketEndpoint(c *fiberWs.Conn) {
	clientsMux.Lock()
	clients[c] = true
	clientsMux.Unlock()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	c.SetPongHandler(func(string) error {
		return c.SetReadDeadline(time.Now().Add(pingInterval + 10*time.Second))
	})

	err := c.SetReadDeadline(time.Now().Add(pingInterval + 10*time.Second))
	if err != nil {
		logrus.Errorf("set deadline error: %v", err)
		return
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			var walletQuery types.BalanceQuery
			err := c.ReadJSON(&walletQuery)
			if err != nil {
				if strings.Contains(err.Error(), "websocket: close 1001") ||
					strings.Contains(err.Error(), "websocket: close 1006") {
					done <- struct{}{}
					return
				}
				logrus.Errorf("read error: %v", err)
				done <- struct{}{}
				return
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

			if err := c.WriteJSON(response); err != nil {
				logrus.Errorf("write error: %v", err)
				done <- struct{}{}
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			goto cleanup
		case <-ticker.C:
			if err := c.WriteMessage(fiberWs.PingMessage, []byte{}); err != nil {
				logrus.Errorf("ping error: %v", err)
				goto cleanup
			}
		}
	}

cleanup:
	clientsMux.Lock()
	delete(clients, c)
	clientsMux.Unlock()
}
