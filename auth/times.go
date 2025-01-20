package auth

import (
	"math"
	"sync"
	"time"
)

type NonceStore struct {
	store map[string]time.Time // nonce -> timestamp
	mu    sync.RWMutex         // Mutex for nonce store
}

type AuthRequest struct {
	Timestamp int64  // Unix timestamp
	Nonce     string // Unique nonce for each request
}

// NonceStoreInstance is a global instance of NonceStore.
var NonceStoreInstance = NewNonceStore()

// NewNonceStore creates a new instance of NonceStore.
// NonceStore is used to store nonces along with their creation time.
// It provides thread-safe access to the stored nonces.
func NewNonceStore() *NonceStore {
	return &NonceStore{
		store: make(map[string]time.Time),
		mu:    sync.RWMutex{},
	}
}

// ValidateRequest validates the given authentication request.
// It checks if the request timestamp is within the allowed time window
// and if the nonce has not been used before.
// If the request is valid, it adds the nonce to the store and returns true.
// Otherwise, it returns false.
func (ns *NonceStore) ValidateRequest(req AuthRequest) bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	currentTime := time.Now().UTC().Unix()
	timeWindow := 15

	if req.Timestamp < 0 ||
		math.Abs(float64(currentTime-req.Timestamp)) > float64(timeWindow) {
		return false
	}

	if _, exists := ns.store[req.Nonce]; exists {
		return false
	}

	ns.store[req.Nonce] = time.Now()
	return true
}
