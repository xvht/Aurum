package types

type BinanceMessage struct {
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
