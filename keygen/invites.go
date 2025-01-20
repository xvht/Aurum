package keygen

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	segmentLength = 10 // Length of each segment
	segmentCount  = 4  // Number of segments

	consonants = "bcdfghjklmnpqrstvwxz" // Consonants
	vowels     = "aeiou"                // Vowels
	numbers    = "0123456789"           // Numbers 0-9
)

// generateRandomIndex generates a random index between 0 and max (exclusive).
func generateRandomIndex(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

// generateSegment generates a random segment based on a predefined pattern.
// It returns a string representing the generated segment.
func generateSegment() string {
	var segment strings.Builder
	pattern := []string{"consonant", "vowel", "number", "consonant"}

	for i := 0; i < segmentLength; i++ {
		var char string
		switch pattern[i%len(pattern)] {
		case "consonant":
			char = string(consonants[generateRandomIndex(len(consonants))])
		case "vowel":
			char = string(vowels[generateRandomIndex(len(vowels))])
		case "number":
			char = string(numbers[generateRandomIndex(len(numbers))])
		}

		if generateRandomIndex(2) == 1 {
			char = strings.ToUpper(char)
		}

		segment.WriteString(char)
	}

	return segment.String()
}

// GenerateInvite generates a random invite code consisting of multiple segments.
// Each segment is generated using the generateSegment function.
// The segments are then joined together using a hyphen (-) as the separator.
// The generated invite code is returned as a string.
func GenerateInvite() string {
	segments := make([]string, segmentCount)
	for i := 0; i < segmentCount; i++ {
		segments[i] = generateSegment()
	}
	return strings.Join(segments, "-")
}

// GenerateBulkInvites generates a specified number of unique invites.
// It returns a slice of strings containing the generated invites.
func GenerateBulkInvites(count int) []string {
	invites := make(map[string]bool)

	for len(invites) < count {
		invite := GenerateInvite()
		invites[invite] = true
	}

	result := make([]string, 0, count)
	for invite := range invites {
		result = append(result, invite)
	}

	return result
}

// ValidateInvite checks if the given invite string is valid.
// It splits the invite string into segments and checks if the number of segments is correct.
// It also checks if each segment has the correct length.
// Returns true if the invite is valid, otherwise returns false.
func PreValidateInvite(invite string) bool {
	segments := strings.Split(invite, "-")
	if len(segments) != segmentCount {
		return false
	}

	for _, segment := range segments {
		if len(segment) != segmentLength {
			return false
		}
	}

	return true
}
