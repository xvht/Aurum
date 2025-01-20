package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultHashCost = 12 // The default cost factor used for password hashing.
)

// PasswordManagerInstance is a global instance of PasswordManager with the default hash cost.
var PasswordManagerInstance = NewPasswordManager()

type PasswordManager struct {
	HashCost int // The cost factor used for password hashing. Higher values increase the time required to hash a password.
}

// NewPasswordManager creates a new instance of PasswordManager with the specified hash cost.
// If no hash cost is provided, it uses the default hash cost.
func NewPasswordManager(cost ...int) *PasswordManager {
	hashCost := DefaultHashCost
	if len(cost) > 0 {
		hashCost = cost[0]
	}
	return &PasswordManager{HashCost: hashCost}
}

// HashPassword hashes the given password using bcrypt algorithm.
// It returns the hashed password as a string and any error encountered during the hashing process.
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.HashCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its corresponding hash and returns true if they match.
// It uses bcrypt.CompareHashAndPassword to perform the comparison.
func (pm *PasswordManager) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
