package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/d0ctr/bldbr-bot-go/shared"
	"github.com/d0ctr/bldbr-bot-go/commands"
)

var logger = slog.Default().With(slog.String("component", "tg-client"))

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
	})

	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		Logger: logger,
	})

	tgClient := &TgClient { bot, dispatcher, updater }

	return tgClient, nil
}

func (tg *TgClient) Start(wg *sync.WaitGroup) {
	for name, def := range commands.All {
		tg.dispatcher.AddHandler(handlers.NewCommand(name, def.Handler))
	}

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
