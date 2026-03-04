package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dotenv-org/godotenvvault"
)

const (
	TELEGRAM_TOKEN = "TELEGRAM_TOKEN"
	ENV = "ENV"
)


func main() {
	logger := NewLogger()
	slog.SetDefault(logger)
	if err := godotenvvault.Load(); err != nil {
		logger.Error("failed to acquire env variables", err);
	}

	tgClient, err := NewTgClient(os.Getenv(TELEGRAM_TOKEN))
	if err != nil {
		logger.Error("failed to start tg bot", err)
	}


	var wg sync.WaitGroup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	tgClient.Start(&wg)
	
	go func() {
		signal := <-stop
		logger.Info("received signal: {}, stopping", signal)

		err := tgClient.Stop()
		if err != nil {
			logger.Error("telegram bot has failed to stop", err)
		}

		if err != nil {
			os.Exit(1)
		}
	}()

	wg.Wait()

	logger.Info("application has finished")
}
