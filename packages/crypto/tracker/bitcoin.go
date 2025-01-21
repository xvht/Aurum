package tracker

import (
	"aurum/packages/crypto/prices"
	"aurum/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type bitcoinBalanceResponse struct {
	FinalBalance int64 `json:"final_balance"`
}

func getBitcoinBalance(address string) (types.WalletData, error) {
	resp, err := http.Get(fmt.Sprintf("https://blockchain.info/rawaddr/%s", address))
	if err != nil {
		return types.WalletData{}, err
	}
	defer resp.Body.Close()

	var data bitcoinBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return types.WalletData{}, err
	}

	balance := float64(data.FinalBalance) / 1e8
	btcPrice, _ := prices.Instance.GetPrice("BTC")

	walletData := types.WalletData{
		TokenSymbol: "BTC",
		Balance:     balance,
		USDValue:    balance * btcPrice.Price,
		LastUpdate:  time.Now().UnixMilli(),
	}

	return walletData, nil
}
