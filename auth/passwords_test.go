package auth_test

import (
	"testing"
	"vexal/auth"

	"github.com/stretchr/testify/assert"
)

func TestNewPasswordManager(t *testing.T) {
	// Test default hash cost
	pm := auth.NewPasswordManager()
	assert.Equal(t, auth.DefaultHashCost, pm.HashCost)

	// Test custom hash cost
	customHashCost := 14
	pm = auth.NewPasswordManager(customHashCost)
	assert.Equal(t, customHashCost, pm.HashCost)
}

func TestHashPassword(t *testing.T) {
	pm := auth.NewPasswordManager()

	// Test hashing a password
	password := "superSecretPassword123"
	hashedPassword, err := pm.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Test hashing an empty password
	emptyPassword := ""
	emptyHashedPassword, err := pm.HashPassword(emptyPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, emptyHashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	pm := auth.NewPasswordManager()

	// Test matching password and hash
	password := "superSecretPassword123"
	hashedPassword, _ := pm.HashPassword(password)
	assert.True(t, pm.CheckPasswordHash(password, hashedPassword))

	// Test non-matching password and hash
	incorrectPassword := "superSecretPassword786"
	assert.False(t, pm.CheckPasswordHash(incorrectPassword, hashedPassword))
}
