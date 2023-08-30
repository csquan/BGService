package db

import (
	"context"
	"fmt"
	"github.com/ethereum/api-in/config"
	"github.com/go-redis/redis/v8"
)

type CustomizedRedis struct {
	*redis.Client
}

var Rdb CustomizedRedis

func GetRedisEngine(config *config.Config) CustomizedRedis {
	Rdb = NewRedis(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
	return Rdb
}

func NewRedis(Addr string, Password string, DB int) CustomizedRedis {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       DB,
	})
	ctx := context.Background()
	pong, err := Rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
	if err != nil {
		fmt.Println(err)
	}
	return CustomizedRedis{Rdb}
}
