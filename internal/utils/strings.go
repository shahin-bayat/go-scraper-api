package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
)

func GenerateRandomString(length int) (string, error) {
	// Create byte slice of appropriate length
	randomBytes := make([]byte, length)

	// Read random bytes from cryptographic RNG
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Encode random bytes to base64 string
	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	// Return the first 'length' characters of the base64 string
	return randomString[:length], nil
}

func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
