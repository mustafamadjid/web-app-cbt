package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateSecretToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func main() {
	token, err := GenerateSecretToken(64) // 32 byte ≈ 256-bit
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}