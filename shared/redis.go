package shared

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func Redis() (*redis.Client, context.Context) {
	return redisClient, context.Background()
}

func init() {
	var logger = slog.Default().With("component", "redis")

	url, ok := REDIS_URL.Get()
	if !ok {
		logger.Error("redis url is empty")
		return
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		logger.Error("failed to parse redis url: {}", err)
	}

	opt.OnConnect = func (context.Context, *redis.Conn) error {
		logger.Info("redis connection established")
		return nil
	}

	redisClient = redis.NewClient(opt)

	redisClient.Ping(context.Background())
}
