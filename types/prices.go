package types

type PriceData struct {
	Symbol     string  `json:"symbol"`
	Price      float64 `json:"price"`
	Change24h  float64 `json:"change24h"`
	Volume     float64 `json:"volume"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	LastUpdate int64   `json:"lastUpdate"`
}
