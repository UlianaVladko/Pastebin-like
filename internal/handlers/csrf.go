package handlers

import (
	"crypto/rand"
	"encoding/hex"
)

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}