package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// Хэширование пароля
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
