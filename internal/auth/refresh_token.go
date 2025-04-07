package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {	//Creates a random 256-bit (32-byte) string key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	tokenString := hex.EncodeToString(key)
	return tokenString, nil
}