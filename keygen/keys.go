package keygen

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"math/big"
	"sort"
	"strings"
	"unicode"
)

const (
	LicenseBytes        = 32  // Number of bytes for license key
	LicenseKeySeparator = '-' // Separator character for license key
)

// GenerateAdvancedLicenseKey generates an advanced license key.
// It returns a string representing the generated license key and an error, if any.
// The function generates random bytes, encodes them using base64 URL-safe encoding,
// removes any non-alphanumeric characters, inserts dashes at random positions,
// and randomizes the case of the alphanumeric characters.
func GenerateAdvancedLicenseKey() (string, error) {
	randomBytes := make([]byte, LicenseBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Printf("Cannot generate random bytes: %v", err)
		return "", err
	}

	// base64 url-safe encode
	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
	alphanumeric := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, encoded)

	separatorPositions := generateSeparatorPositions(alphanumeric)
	encodedWithRandomCase := randomizeCase(alphanumeric)
	result := insertDashes(encodedWithRandomCase, separatorPositions)

	return result, nil
}

func generateSeparatorPositions(s string) []int {
	maxSeparators := len(s) / 6 // reduce separator frequency
	separatorCount, _ := rand.Int(rand.Reader, big.NewInt(int64(maxSeparators-1)+1))

	positions := make([]int, separatorCount.Int64()+1)
	for i := range positions {
		pos, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
		positions[i] = int(pos.Int64())
	}

	// avoid overlapping positions
	sort.Ints(positions)
	return positions
}

func randomizeCase(s string) string {
	result := []rune(s)
	for i := range result {
		shouldUppercase, _ := rand.Int(rand.Reader, big.NewInt(2))
		if shouldUppercase.Int64() == 1 {
			result[i] = unicode.ToUpper(result[i])
		}
	}
	return string(result)
}

func insertDashes(s string, positions []int) string {
	result := []rune(s)

	for i := len(positions) - 1; i >= 0; i-- {
		result = append(result[:positions[i]], append([]rune{LicenseKeySeparator}, result[positions[i]:]...)...)
	}

	return string(result)
}
