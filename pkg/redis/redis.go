package redis

import (
	"context"
	"user/pkg/conf"
	"user/pkg/log"

	"github.com/go-redis/redis/v8"
)

var (
	Cli    *redis.Client
	logger = log.C.Logger().With().Str("base", "redis").Logger()
)

// MustConnectRedis 初始化redis链接，连不上就死
func MustConnectRedis(redisConf *conf.Config) {

	Cli = redis.NewClient(&redis.Options{
		Addr:       redisConf.RedisAddr,
		Password:   redisConf.RedisPass,
		DB:         redisConf.RedisDB,
		MaxRetries: 3,
	})

	ctx := context.Background()
	_, err := Cli.Ping(ctx).Result()

	if err != nil {
		logger.Error().Err(err).Send()
	}

}
