package telegram

import (
	"errors"
	"os"

	l "github.com/d0ctr/bldbr-bot-go/common/logger"

	tele "gopkg.in/telebot.v3"
)

var (
	logger = l.CreateLogger(&l.LoggerOptions{Module: "telegram-bot"})

	TOKEN = os.Getenv("TELEGRAM_TOKEN")

	ENV = os.Getenv("ENV")
)

type Bot tele.Bot

func CreateBot() (*Bot, error) {
	if len(TOKEN) == 0 {
		return nil, errors.New("missing token for starting telegram bot")
	}

	if ENV == "dev" {}

	bot, err := tele.NewBot(tele.Settings{
		Token: TOKEN,
		URL: ,
	})

	return (*Bot)(bot), nil
}
