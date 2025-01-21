package crypto

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberWs "github.com/gofiber/websocket/v2"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type PriceData struct {
	Symbol     string  `json:"symbol"`
	Price      float64 `json:"price"`
	Change24h  float64 `json:"change24h"`
	Volume     float64 `json:"volume"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	LastUpdate int64   `json:"lastUpdate"`
}

type binanceMessage struct {
	EventType    string `json:"e"`
	EventTime    int64  `json:"E"`
	Symbol       string `json:"s"`
	Change24h    string `json:"p"`
	PriceChgPct  string `json:"P"`
	WeightedAvg  string `json:"w"`
	PrevClose    string `json:"x"`
	Price        string `json:"c"`
	LastQty      string `json:"Q"`
	BidPrice     string `json:"b"`
	BidQty       string `json:"B"`
	AskPrice     string `json:"a"`
	AskQty       string `json:"A"`
	Open         string `json:"o"`
	High         string `json:"h"`
	Low          string `json:"l"`
	Volume       string `json:"v"`
	QuoteVolume  string `json:"q"`
	OpenTime     int64  `json:"O"`
	CloseTime    int64  `json:"C"`
	FirstTradeId int64  `json:"F"`
	LastTradeId  int64  `json:"L"`
	TradeCount   int64  `json:"n"`
}

var (
	clients    = make(map[*fiberWs.Conn]bool)
	clientsMux sync.RWMutex
	prices     = make(map[string]PriceData)
	pricesMux  sync.RWMutex
)

func connectToBinance() {
	pairs := []string{"btcusdt", "ethusdt", "ltcusdt", "solusdt"}
	streams := make([]string, len(pairs))
	for i, pair := range pairs {
		streams[i] = pair + "@ticker"
	}

	wsURL := "wss://stream.binance.com:9443/ws/" + strings.Join(streams, "/")

	for {
		c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			logrus.Errorf("binance dial error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		logrus.Infoln("Connected to Binance WebSocket")

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logrus.Errorf("binance read error: %v", err)
				break
			}

			var binanceData binanceMessage
			if err := json.Unmarshal(message, &binanceData); err != nil {
				logrus.Errorf("binance unmarshal error: %v", err)
				continue
			}

			priceData := PriceData{
				Symbol:     binanceData.Symbol,
				Price:      mustParseFloat(binanceData.Price),
				Change24h:  mustParseFloat(binanceData.Change24h),
				Volume:     mustParseFloat(binanceData.Volume),
				High:       mustParseFloat(binanceData.High),
				Low:        mustParseFloat(binanceData.Low),
				LastUpdate: time.Now().UnixMilli(),
			}

			pricesMux.Lock()
			prices[priceData.Symbol] = priceData
			pricesMux.Unlock()

			broadcastPrice(priceData)
		}

		c.Close()
		logrus.Warnln("Binance WebSocket disconnected, reconnecting...")
		time.Sleep(5 * time.Second)
	}
}

func broadcastPrice(data PriceData) {
	clientsMux.RLock()
	defer clientsMux.RUnlock()

	for client := range clients {
		if err := client.WriteJSON(data); err != nil {
			logrus.Errorf("write error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func mustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

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
	pricesMux.RLock()
	for _, price := range prices {
		err := c.WriteJSON(price)
		if err != nil {
			logrus.Errorf("write error: %v", err)
			break
		}
	}
	pricesMux.RUnlock()

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
	go connectToBinance()
}
