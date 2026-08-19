package services

import (
	"context"
	"d0ctr/bldbr-bot/shared"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func Redis() (*redis.Client, context.Context) {
	return redisClient, context.Background()
}

func init() {
	var logger = shared.WithComponent("redis")

	url, ok := shared.REDIS_URL.Get()
	if !ok {
		logger.Error("redis url is empty")
		return
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		logger.Error("failed to parse redis url: {}", shared.ErrAttr(err))
	}

	opt.OnConnect = func (context.Context, *redis.Conn) error {
		logger.Info("redis connection established")
		return nil
	}

	redisClient = redis.NewClient(opt)

	redisClient.Ping(context.Background())
}
