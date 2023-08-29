package util

import (
	"math/rand"
	"time"
)

func GenerateVerifyCode(length int) string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}
