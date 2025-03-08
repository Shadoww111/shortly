package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortCode creates a random alphanumeric code.
func GenerateShortCode(length int) (string, error) {
	var sb strings.Builder
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		sb.WriteByte(charset[n.Int64()])
	}
	return sb.String(), nil
}

// IsValidCustomCode checks if a custom short code is valid.
func IsValidCustomCode(code string) bool {
	if len(code) < 3 || len(code) > 20 {
		return false
	}
	for _, c := range code {
		if !strings.ContainsRune(charset+"-_", c) {
			return false
		}
	}
	return true
}
