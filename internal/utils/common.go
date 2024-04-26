package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func GenerateRandomString(length uint) (string, error) {
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

//func StringInSlice(slice []string, str string) bool {
//	for _, s := range slice {
//		if s == TrimSpaceLower(str) {
//			return true
//		}
//	}
//	return false
//}

func KeyInMap(m map[string]string, key string) bool {
	_, ok := m[key]
	return ok
}

func TrimSpaceLower(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}

func UintInSlice(slice []uint, value uint) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
