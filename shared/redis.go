package shared

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func Redis() (*redis.Client, context.Context) {
	if (redisClient == nil) {
		redisClient = getRedisClient()
	}
	return redisClient, context.Background()
}

var redisClient *redis.Client = nil 

func getRedisClient() *redis.Client {
	var logger = slog.Default().With("component", "redis")

	url, ok := REDIS_URL.Get()
	if !ok {
		logger.Warn("redis url is empty")
		return nil
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		logger.Error("failed to parse redis url", err)
		return nil
	}

	opt.OnConnect = func (context.Context, *redis.Conn) error {
		logger.Info("redis connection established")
		return nil
	}

	return redis.NewClient(opt)
}
