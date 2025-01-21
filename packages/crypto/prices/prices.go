package prices

import (
	"aurum/types"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type PriceTracker struct {
	prices    map[string]types.PriceData
	pricesMux sync.RWMutex
	pairs     []string
}

var Instance = NewPriceTracker()

func NewPriceTracker() *PriceTracker {
	tracker := &PriceTracker{
		prices: make(map[string]types.PriceData),
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

			var binanceData types.BinanceMessage
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

func (pt *PriceTracker) updatePrice(data types.BinanceMessage) {
	pt.pricesMux.Lock()
	defer pt.pricesMux.Unlock()

	pt.prices[data.Symbol] = types.PriceData{
		Symbol:     data.Symbol,
		Price:      mustParseFloat(data.Price),
		Change24h:  mustParseFloat(data.Change24h),
		Volume:     mustParseFloat(data.Volume),
		High:       mustParseFloat(data.High),
		Low:        mustParseFloat(data.Low),
		LastUpdate: time.Now().UnixMilli(),
	}
}

func (pt *PriceTracker) GetPrice(symbol string) (types.PriceData, bool) {
	pt.pricesMux.RLock()
	defer pt.pricesMux.RUnlock()

	price, exists := pt.prices[fmt.Sprintf("%s%s", symbol, "USDT")]
	return price, exists
}

func (pt *PriceTracker) GetAllPrices() map[string]types.PriceData {
	pt.pricesMux.RLock()
	defer pt.pricesMux.RUnlock()

	prices := make(map[string]types.PriceData, len(pt.prices))
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
