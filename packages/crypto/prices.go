package crypto

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

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

type PriceTracker struct {
	prices    map[string]PriceData
	pricesMux sync.RWMutex
	pairs     []string
}

func NewPriceTracker() *PriceTracker {
	tracker := &PriceTracker{
		prices: make(map[string]PriceData),
		pairs:  []string{"btcusdt", "ethusdt", "ltcusdt", "solusdt"},
	}
	go tracker.connectToBinance()
	return tracker
}

func (pt *PriceTracker) connectToBinance() {
	streams := make([]string, len(pt.pairs))
	for i, pair := range pt.pairs {
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

			pt.updatePrice(binanceData)
		}

		c.Close()
		logrus.Warnln("Binance WebSocket disconnected, reconnecting...")
		time.Sleep(5 * time.Second)
	}
}

func (pt *PriceTracker) updatePrice(data binanceMessage) {
	pt.pricesMux.Lock()
	defer pt.pricesMux.Unlock()

	pt.prices[data.Symbol] = PriceData{
		Symbol:     data.Symbol,
		Price:      mustParseFloat(data.Price),
		Change24h:  mustParseFloat(data.Change24h),
		Volume:     mustParseFloat(data.Volume),
		High:       mustParseFloat(data.High),
		Low:        mustParseFloat(data.Low),
		LastUpdate: time.Now().UnixMilli(),
	}
}

func (pt *PriceTracker) GetPrice(symbol string) (PriceData, bool) {
	pt.pricesMux.RLock()
	defer pt.pricesMux.RUnlock()

	price, exists := pt.prices[symbol]
	return price, exists
}

func (pt *PriceTracker) GetAllPrices() map[string]PriceData {
	pt.pricesMux.RLock()
	defer pt.pricesMux.RUnlock()

	// Return a copy to prevent external modifications
	prices := make(map[string]PriceData, len(pt.prices))
	for k, v := range pt.prices {
		prices[k] = v
	}
	return prices
}

func mustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
