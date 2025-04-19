package bft

// defense.go: Cryptographic defensive mechanisms against attack vectors
// Implements cryptographic defenses for BFT.

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HMAC generates a keyed-hash message authentication code using SHA-256.
func HMAC(key, message []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyHMAC verifies the HMAC of a message.
func VerifyHMAC(key, message []byte, messageMAC string) bool {
	expectedMAC := HMAC(key, message)
	return hmac.Equal([]byte(expectedMAC), []byte(messageMAC))
}
