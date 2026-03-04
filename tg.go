package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/d0ctr/bldbr-bot-go/commands"
)

type TgClient struct {
	bot *gotgbot.Bot
	dispatcher *ext.Dispatcher
	updater *ext.Updater
	logger *slog.Logger
}

func NewTgClient(token string) (*TgClient, error) {
	logger := slog.Default().With(slog.String("component", "tg-client"))
	if len(token) == 0 {
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

	tgClient := &TgClient { bot, dispatcher, updater, logger }

	return tgClient, nil
}

func (tg *TgClient) Start(wg *sync.WaitGroup) {
	tg.dispatcher.AddHandler(handlers.NewCommand("ping", commands.Ping))
	tg.dispatcher.AddHandler(handlers.NewCommand("ahegao", commands.Ahegao))

	tg.updater.StartPolling(tg.bot, nil)
	tg.logger.Info("telegram bot has started")

	wg.Go(func() {
		tg.updater.Idle()
		tg.logger.Info("telegram bot has finished")
	})
}

func (tg *TgClient) Stop() error {
	return tg.updater.Stop()
}
