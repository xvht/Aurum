package tracker

import (
	"aurum/packages/crypto/prices"
	"aurum/types"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RPCResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Result  map[string]interface{} `json:"result"`
}

func getSolanaBalance(address string) (types.WalletData, error) {
	rpcRequest := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getBalance",
		Params:  []interface{}{address},
	}

	requestBody, err := json.Marshal(rpcRequest)
	if err != nil {
		return types.WalletData{}, err
	}

	resp, err := http.Post("https://api.mainnet-beta.solana.com", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return types.WalletData{}, err
	}

	defer resp.Body.Close()

	var rpcResponse RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return types.WalletData{}, err
	}

	balanceLamports := rpcResponse.Result["value"].(float64)
	balance := balanceLamports / 1e9
	solPrice, _ := prices.Instance.GetPrice("SOL")

	walletData := types.WalletData{
		TokenSymbol: "SOL",
		Balance:     balance,
		USDValue:    balance * solPrice.Price,
		LastUpdate:  time.Now().UnixMilli(),
	}

	return walletData, nil
}
