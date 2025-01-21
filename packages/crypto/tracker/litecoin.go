package tracker

import (
	"aurum/packages/crypto/prices"
	"aurum/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type data struct {
	Address struct {
		Balance float64 `json:"balance"`
	} `json:"address"`
}

type litecoinBalanceResponse struct {
	Data map[string]data `json:"data"`
}

func getLitecoinBalance(address string) (types.WalletData, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.blockchair.com/litecoin/dashboards/address/%s?limit=0", address))
	if err != nil {
		return types.WalletData{}, err
	}
	defer resp.Body.Close()

	var data litecoinBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return types.WalletData{}, err
	}

	balance := float64(data.Data[address].Address.Balance) / 1e8
	ltcPrice, _ := prices.Instance.GetPrice("LTC")

	walletData := types.WalletData{
		TokenSymbol: "LTC",
		Balance:     balance,
		USDValue:    balance * ltcPrice.Price,
		LastUpdate:  time.Now().UnixMilli(),
	}

	return walletData, nil
}
