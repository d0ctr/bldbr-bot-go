package commands

import (
	"github.com/d0ctr/bldbr-bot-go/common/locale"
	"github.com/d0ctr/bldbr-bot-go/common/logger"
	"gopkg.in/telebot.v3"
)

type command struct {
	name string
}

type Command interface {
	Name() string
	Description(locale.Language) string
}

func (c *command) Name() string {
	return c.name
}

func (c *command) Description(lang locale.Language) string {
	line, err := locale.Get(lang, "command_"+c.name+"_description")
	if err != nil {
		panic(err)
	}
	return line
}

type TelegramContext = telebot.Context

type TelegramCommand interface {
	Telegram(TelegramContext) error
}

func getLoggerFromContext(ctx TelegramContext) func(string) *logger.Logger {
	switch l := ctx.Get("Logger").(type) {
	case func(string) *logger.Logger:
		return l
	default:
		return func(s string) *logger.Logger { return logger.RootLogger.Child(&logger.ChildLoggerOptions{Module: s}) }
	}
}

func (c *command) telegramLog(ctx TelegramContext) *logger.Logger {
	l := getLoggerFromContext(ctx)(c.name)

	l.F.DEBUG("received command [%s] with args [%v]", c.name, ctx.Args())

	return l
}
