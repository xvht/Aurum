package tracker

import (
	"aurum/types"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	WalletCache map[string]types.WalletData
	cacheMutex  sync.Mutex
	cacheExpiry = 5 * time.Minute
)

func init() {
	WalletCache = make(map[string]types.WalletData)
}

func GetWalletBalance(walletType, walletAddress string) (types.WalletData, bool, error) {
	cacheKey := getWalletCacheKey(walletType, walletAddress)
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if balance, ok := WalletCache[cacheKey]; ok && time.Now().Sub(time.UnixMilli(balance.LastUpdate)) < cacheExpiry {
		return balance, true, nil
	}

	balance, err := fetchWalletBalance(walletType, walletAddress)
	if err != nil {
		logrus.Errorf("failed to fetch wallet balance: %v", err)
		return types.WalletData{}, false, err
	}

	balance.LastUpdate = time.Now().UnixMilli()
	WalletCache[cacheKey] = balance
	return balance, true, nil
}

func fetchWalletBalance(walletType, walletAddress string) (types.WalletData, error) {
	switch walletType {
	case "ETH":
		return getEthereumBalance(walletAddress)
	case "BTC":
		return getBitcoinBalance(walletAddress)
	case "LTC":
		return getLitecoinBalance(walletAddress)
	case "SOL":
		return getSolanaBalance(walletAddress)
	default:
		return types.WalletData{}, nil
	}
}

func getWalletCacheKey(walletType, walletAddress string) string {
	return walletType + "-" + walletAddress
}
