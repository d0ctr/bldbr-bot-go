package tg

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"d0ctr/bldbr-bot/shared"
	"d0ctr/bldbr-bot/tg/commands"
	"d0ctr/bldbr-bot/tg/handlers"
	"d0ctr/bldbr-bot/tg/utils"
)

var logger = slog.Default().With("component", "tg-client")

type TgClient struct {
	bot *gotgbot.Bot
	dispatcher *ext.Dispatcher
	updater *ext.Updater
}

func NewTgClient() (*TgClient, error) {
	token, ok := shared.TELEGRAM_TOKEN.Get()
	if !ok {
		return nil, fmt.Errorf("token is empty")
	}

	opts := &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			UseTestEnvironment: false, // env == "local" || env == "dev",
		},
	}

	bot, err := gotgbot.NewBot(token, opts)
	if err != nil {
		return nil, err
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Logger: logger,
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			text := fmt.Sprintf("Ошибка:\n<pre>%v</pre>", err)

			if _, ok := err.(utils.NoSendError); !ok {
				if ctx.EffectiveMessage != nil {
					ctx.EffectiveMessage.Reply(b, text, utils.GetDefaultMessageOpts())
				} else if ctx.EffectiveChat != nil {
					ctx.EffectiveChat.SendMessage(b, text, utils.GetDefaultMessageOpts())
				}
			}

			slog.Error("got an error", err)

			return ext.DispatcherActionNoop
		},
	})

	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		Logger: logger,
	})

	tgClient := &TgClient { bot, dispatcher, updater }

	return tgClient, nil
}

func (tg *TgClient) Start(wg *sync.WaitGroup) {
	tg.dispatcher.AddHandler(handlers.AutoDelete())
	tg.dispatcher.AddHandler(handlers.ReplyToBot())
	for handler := range handlers.Commands() {
		tg.dispatcher.AddHandler(handler)
	}

	tg.bot.DeleteWebhook(&gotgbot.DeleteWebhookOpts{ DropPendingUpdates: false })

	tg.updater.StartPolling(tg.bot, nil)
	logger.Info("telegram bot has started")

	tg.PublishCommands()

	wg.Go(func() {
		tg.updater.Idle()
		logger.Info("telegram bot has finished")
	})
}

func (tg *TgClient) Stop() error {
	return tg.updater.Stop()
}

func (tg *TgClient) PublishCommands() {
	var botCommands []gotgbot.BotCommand
	for _, def := range commands.All {
		botCommand := gotgbot.BotCommand{
			Command: def.Name,
			Description: def.Description,
		}

		botCommands = append(botCommands, botCommand)
	}

	if ok, err := tg.bot.SetMyCommands(botCommands, nil); !ok {
		logger.Warn("failed to set commands", err)
	} else {
		logger.Debug("have set commands")
	}

}

