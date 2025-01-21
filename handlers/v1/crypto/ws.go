package crypto

import (
	"aurum/packages/crypto"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberWs "github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
)

var (
	clients    = make(map[*fiberWs.Conn]bool)
	clientsMux sync.RWMutex
	tracker    *crypto.PriceTracker
)

func WebSocketHandler(c *fiber.Ctx) error {
	if fiberWs.IsWebSocketUpgrade(c) {
		logrus.Infof("WebSocket connection from %s", c.IP())
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
	priceData := tracker.GetAllPrices()
	for _, price := range priceData {
		err := c.WriteJSON(price)
		if err != nil {
			logrus.Errorf("write error: %v", err)
			break
		}
	}

	go func() {
		for {
			priceData := tracker.GetAllPrices()
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

	// Unregister client on disconnect
	clientsMux.Lock()
	delete(clients, c)
	clientsMux.Unlock()
}

func init() {
	tracker = crypto.NewPriceTracker()
}
