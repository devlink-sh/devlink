package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateUUID generates a random UUID
func GenerateUUID() string {
	// Generate 16 random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("id-%d", time.Now().UnixNano())
	}

	// Set version (4) and variant bits
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	// Convert to hex string
	return hex.EncodeToString(bytes[:4]) + "-" +
		hex.EncodeToString(bytes[4:6]) + "-" +
		hex.EncodeToString(bytes[6:8]) + "-" +
		hex.EncodeToString(bytes[8:10]) + "-" +
		hex.EncodeToString(bytes[10:])
}
