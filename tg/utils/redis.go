package utils

import (
	"context"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"d0ctr/bldbr-bot/shared"
	gredis "github.com/redis/go-redis/v9"
)

var redis *gredis.Client

func init() {
	redis, _ = shared.Redis()
}


func GetChatValue(ctx *ext.Context, name string) (string, error) {
	key := fmt.Sprintf("tg:chat:%v", ctx.EffectiveChat.Id)

	if exists, err := redis.HExists(context.Background(), key, name).Result(); !exists || err != nil {
		return "", err
	}

	return redis.HGet(context.Background(), key, name).Result()
}

func SetChatValue(ctx *ext.Context, name string, value string) error {
	key := fmt.Sprintf("tg:chat:%v", ctx.EffectiveChat.Id)

	_, err := redis.HSet(context.Background(), key, name, value).Result()

	return err
}
