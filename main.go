package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"d0ctr/bldbr-bot/shared"
	_ "d0ctr/bldbr-bot/shared"
	"d0ctr/bldbr-bot/tg"
)

func main() {
	tgClient, err := tg.NewTgClient()

	if err != nil {
		slog.Error("failed to start tg bot", shared.ErrAttr(err))
		return
	}

	var wg sync.WaitGroup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	tgClient.Start(&wg)
	
	go func() {
		signal := <-stop
		slog.Info("received signal: {}, stopping", shared.TemplateAttr(signal))

		err := tgClient.Stop()
		if err != nil {
			slog.Error("telegram bot has failed to stop", shared.ErrAttr(err))
		}

		if err != nil {
			os.Exit(1)
		}
	}()

	wg.Wait()

	slog.Info("application has finished")
}
