package types

type BalanceQuery struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`
}

type BalanceResponse struct {
	Error bool       `json:"error"`
	Code  int        `json:"code"`
	Data  WalletData `json:"data"`
}
