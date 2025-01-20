package auth

import (
	"encoding/base64"
	"strings"
)

func Encode(input string) string {
	// Standard encoding
	encoded := base64.StdEncoding.EncodeToString([]byte(input))

	// URL-safe modifications
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.TrimRight(encoded, "=")

	return encoded
}

func Decode(input string) (string, error) {
	// Restore standard base64 padding
	padded := input
	switch len(input) % 4 {
	case 2:
		padded += "=="
	case 3:
		padded += "="
	}

	// Restore standard base64 chars
	padded = strings.ReplaceAll(padded, "-", "+")
	padded = strings.ReplaceAll(padded, "_", "/")

	// Decode
	decoded, err := base64.StdEncoding.DecodeString(padded)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
