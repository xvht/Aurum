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
	go tracker.connectToBitfinex()
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
				logrus.Errorf("Binance read error: %v", err)
				break
			}

			var binanceData types.BinanceMessage
			if err := json.Unmarshal(message, &binanceData); err != nil {
				logrus.Errorf("Binance unmarshal error: %v", err)
				continue
			}

			pt.updatePrice(binanceData)
			logrus.Debugf("Received price update: %s %s", binanceData.Symbol, binanceData.Price)
		}

		c.Close()
		logrus.Warnln("Binance WebSocket disconnected, reconnecting...")
		time.Sleep(5 * time.Second)
	}
}

func (pt *PriceTracker) connectToBitfinex() {
	wsURL := "wss://api-pub.bitfinex.com/ws/2"

	for {
		c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			logrus.Errorf("bitfinex dial error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		logrus.Infoln("Connected to Bitfinex WebSocket")

		subscribeMsg := map[string]interface{}{
			"event":   "subscribe",
			"channel": "ticker",
			"symbol":  "tXMRUSD",
		}

		if err := c.WriteJSON(subscribeMsg); err != nil {
			logrus.Errorf("Failed to subscribe to Bitfinex: %v", err)
			c.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logrus.Errorf("Bitfinex read error: %v", err)
				break
			}

			var data []interface{}
			if err := json.Unmarshal(message, &data); err != nil {
				continue
			}

			pt.processBitfinexData(data)
		}

		c.Close()
		logrus.Warnln("Bitfinex WebSocket disconnected, reconnecting...")
		time.Sleep(5 * time.Second)
	}
}

func (pt *PriceTracker) processBitfinexData(data []interface{}) {
	if !isValidBitfinexTickerData(data) {
		return
	}

	// Extract ticker data array (data[1])
	// https://docs.bitfinex.com/reference/ws-public-ticker#stream-data
	tickerData, ok := data[1].([]interface{})
	if !ok {
		logrus.Debugln("Invalid Bitfinex ticker format")
		return
	}

	// Invalid ticker data
	if len(tickerData) < 10 {
		return
	}

	// LAST_PRICE is at index 6
	// https://docs.bitfinex.com/reference/ws-public-ticker#trading-snapshotupdate-index-1
	lastPrice, ok := tickerData[6].(float64)
	if !ok {
		lastPriceStr, isString := tickerData[6].(string)
		if isString {
			var err error
			lastPrice, err = strconv.ParseFloat(lastPriceStr, 64)
			if err != nil {
				logrus.Errorf("Failed to parse Bitfinex price: %v", err)
				return
			}
		} else {
			logrus.Errorf("Invalid Bitfinex price format")
			return
		}
	}

	dailyChange, _ := extractFloat(tickerData, 4)
	volume, _ := extractFloat(tickerData, 7)
	high, _ := extractFloat(tickerData, 8)
	low, _ := extractFloat(tickerData, 9)

	pt.pricesMux.Lock()
	defer pt.pricesMux.Unlock()

	pt.prices["XMRUSDT"] = types.PriceData{
		Symbol:     "XMRUSDT",
		Price:      lastPrice,
		Change24h:  dailyChange,
		Volume:     volume,
		High:       high,
		Low:        low,
		LastUpdate: time.Now().UnixMilli(),
	}

	logrus.Debugf("Received price update: XMRUSDT %f", lastPrice)
}

func extractFloat(data []interface{}, index int) (float64, bool) {
	if len(data) <= index {
		return 0, false
	}

	switch v := data[index].(type) {
	case float64:
		return v, true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

func isValidBitfinexTickerData(data []interface{}) bool {
	// An array with at least 2 elements
	if len(data) < 2 {
		return false
	}

	// First element is a number (channel ID)
	_, isNum := data[0].(float64)
	if !isNum {
		return false
	}

	// Second element is an array (ticker data) and not a heartbeat
	_, isArray := data[1].([]interface{})
	if !isArray {
		strData, isString := data[1].(string)
		return !isString || strData != "hb" // Not a heartbeat
	}

	return true
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
