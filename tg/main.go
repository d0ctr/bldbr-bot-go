package tg

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"

	"d0ctr/bldbr-bot/shared"
	"d0ctr/bldbr-bot/tg/commands"
	"d0ctr/bldbr-bot/tg/handlers"
	"d0ctr/bldbr-bot/tg/utils"
)

var logger = slog.Default().With(shared.ComponentAttr("tg-client"))

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

			slog.Error("got an error", shared.ErrAttr(err))

			return ext.DispatcherActionNoop
		},
	})

	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		Logger: logger,
	})

	tgClient := &TgClient { bot, dispatcher, updater }

	return tgClient, nil
}

func (tg TgClient) Start(wg *sync.WaitGroup) {
	tg.dispatcher.AddHandler(handlers.AutoDelete())
	tg.dispatcher.AddHandler(handlers.ReplyToBot())
	for handler := range handlers.Commands() {
		tg.dispatcher.AddHandler(handler)
	}

	if err := tg.startWebhook(); err != nil {
		logger.Error("failed to start webhook", shared.ErrAttr(err))

		tg.startPolling()
	}

	tg.PublishCommands()

	wg.Go(func() {
		tg.updater.Idle()
		logger.Info("telegram bot has finished")
	})
}

func (tg TgClient) startPolling() {
	err := tg.updater.StartPolling(tg.bot, &ext.PollingOpts{ EnableWebhookDeletion: true })
	if err != nil {
		log.Panicf("failed to start polling: %v", err)
	}
	logger.Info("telegram bot has started polling")
}

func (tg TgClient) startWebhook() error {
	domain, ok := shared.DOMAIN_URL.Get()
	if !ok {
		return fmt.Errorf("'%v' is not set", shared.DOMAIN_URL)
	}

	secretToken := uuid.NewString()
	urlPath := "webhook/" + uuid.NewString()

	if err := tg.updater.StartWebhook(tg.bot, urlPath, ext.WebhookOpts{ ListenAddr: ":8080", SecretToken: secretToken }); err != nil {
		return fmt.Errorf("failed to start webhook: %v", err)
	}

	url := domain + "/" + urlPath
	if ok, err := tg.bot.SetWebhook(url, &gotgbot.SetWebhookOpts{ SecretToken: secretToken }); !ok {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	logger.Info("telegram bot has started webhook")

	return nil
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
		logger.Warn("failed to set commands", shared.ErrAttr(err))
	} else {
		logger.Debug("have set commands")
	}
}

