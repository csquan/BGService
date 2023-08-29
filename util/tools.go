package util

import (
	"github.com/ethereum/api-in/types"
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

func ResponseMsg(Code int, Message string, Data interface{}) types.HttpRes {
	res := types.HttpRes{}

	res.Code = Code
	res.Message = Message
	res.Data = Data
	return res
}
