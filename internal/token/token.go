package token

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
)

func GenerateToken() string {
	// Step 1: Generate random float64
	r := rand.Float64()

	// Step 2: Convert to string
	rStr := strconv.FormatFloat(r, 'f', -1, 64)

	// Step 3: Hash the string with SHA-256
	hash := sha256.Sum256([]byte(rStr))

	// Step 4: Encode to hex
	hashHex := hex.EncodeToString(hash[:])

	return hashHex
}
