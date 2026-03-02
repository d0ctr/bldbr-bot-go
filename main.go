package main

import (
	"os"
	"log/slog"
	"net/http"

	"github.com/dotenv-org/godotenvvault"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

const (
	TELEGRAM_TOKEN = "TELEGRAM_TOKEN"
	ENV = "ENV"
)

func createLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if (a.Key == slog.TimeKey) {
				return slog.Time(a.Key, a.Value.Time().UTC())
			}

			return a
		},
	};
	handler := slog.NewJSONHandler(os.Stdout, opts);
	logger := slog.New(handler);
	slog.SetDefault(logger);
	return logger;
}

func createTgBot() (*gotgbot.Bot, *ext.Dispatcher, *ext.Updater, error) {
	tgToken := os.Getenv(TELEGRAM_TOKEN)
	// env := os.Getenv(ENV)

	opts := &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			UseTestEnvironment: false, // env == "local" || env == "dev",
		},
	}

	bot, err := gotgbot.NewBot(tgToken, opts)
	if err != nil {
		return nil, nil, nil, err
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Logger: slog.Default(),
	})

	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		Logger: slog.Default(),
	})

	return bot, dispatcher, updater, nil
}

func main() {
	logger := createLogger()
	if err := godotenvvault.Load(); err != nil {
		logger.Error("failed to acquire env variables", slog.Any("error", err));
	}

	bot, dispatcher, updater, err := createTgBot()
	if err != nil {
		logger.Error("failed to start tg bot", slog.Any("error", err))
	}

	dispatcher.AddHandler(handlers.NewCommand(
		"ping",
		func(bot *gotgbot.Bot, ctx *ext.Context) error {
			ctx.EffectiveMessage.Reply(bot, "<code>pong</code>", &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			})
			return nil
		},
	))

	updater.StartPolling(bot, nil)

	updater.Idle()
	
	logger.Info("bot has finished")
}
