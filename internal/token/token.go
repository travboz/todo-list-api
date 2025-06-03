package token

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken generates a cryptographically secure random token of the given byte length.
// The token is returned as a hex-encoded string.
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
