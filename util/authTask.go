package util

import (
	"github.com/ethereum/api-in/db"
	"github.com/ethereum/api-in/types"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

func GenerateCode(length int) string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

func GenerateInviteCode(length int) string {
	// 生成固定长度的随机邀请码
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机种子
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
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

func CheckVerifyCode(c *gin.Context, a db.CustomizedRedis, key string, value string) bool {
	storedCode, err := a.Get(c, key).Result()
	if err == redis.Nil {
		// Key does not exist in Redis
		return false
	} else if err != nil {
		// Error while getting value from Redis
		return false
	}
	return storedCode == value
}
