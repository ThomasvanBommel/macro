package util

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomBytes generates a slice of random bytes of the specified length.
func GenerateRandomBytes(n int) []byte {
	defer Trace("GenerateRandomBytes(n)", "n", n)()

	b := make([]byte, n)
	rand.Read(b) // Documented to never return an error on all but legacy Linux systems

	return b
}

// GenerateRandomHexString generates a random hexadecimal string of the specified length (in bytes).
func GenerateRandomHexString(n int) string {
	defer Trace("GenerateRandomHexString(n)", "n", n)()

	b := GenerateRandomBytes(n)
	return hex.EncodeToString(b)
}
