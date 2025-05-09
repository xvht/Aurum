package prices

import (
	"aurum/packages/crypto/prices"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberWs "github.com/gofiber/websocket/v2"
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
	// Register new client
	clientsMux.Lock()
	clients[c] = true
	clientsMux.Unlock()

	// Send current prices immediately
	priceData := prices.Instance.GetAllPrices()
	for _, price := range priceData {
		err := c.WriteJSON(price)
		if err != nil {
			logrus.Errorf("write error: %v", err)
			break
		}
	}

	go func() {
		for {
			priceData := prices.Instance.GetAllPrices()
			for _, price := range priceData {
				clientsMux.RLock()
				for client := range clients {
					err := client.WriteJSON(price)
					if err != nil {
						logrus.Errorf("write error: %v", err)
						break
					}
				}
				clientsMux.RUnlock()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Keep alive
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}

	clientsMux.Lock()
	delete(clients, c)
	clientsMux.Unlock()
}
