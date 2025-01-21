package types

type WalletData struct {
	Balance     float64 `json:"balance"`
	USDValue    float64 `json:"usdValue"`
	TokenSymbol string  `json:"tokenSymbol"`
	LastUpdate  int64   `json:"lastUpdated"`
}
