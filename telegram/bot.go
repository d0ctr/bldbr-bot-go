package telegram

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/d0ctr/bldbr-bot-go/common/commands"
	"github.com/d0ctr/bldbr-bot-go/common/locale"
	"github.com/d0ctr/bldbr-bot-go/common/logger"

	tele "gopkg.in/telebot.v3"
)

type Bot = tele.Bot
type HandlerFunc = tele.HandlerFunc
type Context = tele.Context

var (
	l = func() *logger.Logger {
		return logger.RootLogger.Child(&logger.ChildLoggerOptions{Module: "telegram-bot"})
	}
	lconfig = func() *logger.Logger {
		return l().Child(&logger.ChildLoggerOptions{Module: "config"})
	}
)

var testURLBuilder = func(URL string, token string, method string) string {
	return URL + "/bot" + token + "/test/" + method
}

var handleError = func(err error, ctx tele.Context) {
	l().ERROR("got error from telegram: ", err)
}

type MyBot struct {
	*tele.Bot
	enabledCommands []interface{}
	Logger          *logger.Logger
}

func CreateBot() (*MyBot, error) {
	TOKEN := os.Getenv("TELEGRAM_TOKEN")
	ENV := os.Getenv("ENV")
	lconfig().F.DEBUG("[%s] is %d characters long", "TELEGRAM_TOKEN", strings.Count(TOKEN, ""))
	lconfig().F.DEBUG("[%s] is %s", "ENV", ENV)

	if strings.Count(TOKEN, "") == 0 {
		return nil, errors.New("missing token for starting telegram bot")
	}
	setting := tele.Settings{
		Token: TOKEN,
		Poller: &tele.LongPoller{
			Timeout:        5 * time.Second,
			AllowedUpdates: []string{"message"},
		},
		OnError:   handleError,
		ParseMode: tele.ModeHTML,
	}
	if strings.EqualFold(ENV, "dev") {
		setting.URLBuilder = testURLBuilder
	}
	b, err := tele.NewBot(setting)
	if err != nil {
		l().ERROR("error starting bot: ", err)
		return nil, err
	}

	bot := &MyBot{
		Bot: b,
		Logger: logger.RootLogger.Child(&logger.ChildLoggerOptions{
			Module: "telegram:" + b.Me.Username,
		}),
	}

	bot.AddMiddleware("Context Logger", func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			ctx.Set("Logger", func(s string) *logger.Logger {
				return bot.Logger.Child(&logger.ChildLoggerOptions{Module: s})
			})
			return next(ctx)
		}
	})

	return bot, nil
}

func (bot *MyBot) AddMiddleware(description string, m tele.MiddlewareFunc) {
	bot.Use(m)
	l().F.DEBUG("added middleware %s", description)
}

func (bot *MyBot) AddHandler(what any, c interface{}, m ...tele.MiddlewareFunc) error {
	command := c.(commands.TelegramCommand)
	bot.Handle(what, command.Telegram, m...)
	bot.enabledCommands = append(bot.enabledCommands, c)
	l().F.DEBUG("will handle [%s]\n", what)
	return nil
}

func (bot *MyBot) setMyCommandsForLang(lang locale.Language) error {
	var commandsArr []tele.Command
	for _, c := range bot.enabledCommands {
		telecommand := tele.Command{
			Text:        c.(commands.Command).Name(),
			Description: c.(commands.Command).Description(lang),
		}
		c = append(commandsArr, telecommand)
	}
	l().F.NOISE("registering [%d] commands in locale[%s] ", len(commandsArr), lang)
	return bot.SetCommands(
		commandsArr,
		lang,
		&tele.CommandScope{Type: tele.CommandScopeDefault},
	)
}

func (bot *MyBot) SetMyCommands() {
	for lang := range locale.AvailableLocales {
		if lang == "default" {
			lang = ""
		}
		err := bot.setMyCommandsForLang(lang)
		if err != nil {
			l().F.WARN("error while setting commands [%d] for locale [%s]", len(bot.enabledCommands), lang)
		}
	}
	l().DEBUG("done registering commands")
}
