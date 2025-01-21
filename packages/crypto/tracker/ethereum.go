package tracker

import (
	"aurum/env"
	"aurum/packages/crypto/prices"
	"aurum/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ethereumBalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func getEthereumBalance(address string) (types.WalletData, error) {
	qParams := fmt.Sprintf("module=account&action=balance&address=%s&tag=latest&apikey=%s", address, env.ETHERSCAN_API_KEY)
	resp, err := http.Get(fmt.Sprintf("https://api.etherscan.io/api?%s", qParams))
	if err != nil {
		return types.WalletData{}, err
	}
	defer resp.Body.Close()

	var data ethereumBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return types.WalletData{}, err
	}

	result, err := strconv.ParseFloat(data.Result, 64)
	if err != nil {
		return types.WalletData{}, err
	}
	balance := result / 1e18
	ethPrice, _ := prices.Instance.GetPrice("ETH")

	walletData := types.WalletData{
		TokenSymbol: "ETH",
		Balance:     balance,
		USDValue:    balance * ethPrice.Price,
		LastUpdate:  time.Now().UnixMilli(),
	}

	return walletData, nil
}
