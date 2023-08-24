package main

import (
	"os"
	"strconv"

	"github.com/d0ctr/bldbr-bot-go/common/commands"
	"github.com/d0ctr/bldbr-bot-go/common/logger"
	"github.com/d0ctr/bldbr-bot-go/telegram"

	"github.com/dotenv-org/godotenvvault"
)

func main() {
	logger.RootLogger.INFO("Starting...")
	err := godotenvvault.Load()
	if err != nil {
		logger.RootLogger.FATAL("error loading dotenv: ", err)
	}

	logger.SetRootLogger(logger.CreateLogger(&logger.LoggerOptions{
		Module:    os.Getenv("ENV"),
		IsColored: func() bool { v, _ := strconv.ParseBool(os.Getenv("COLORED_LOG")); return v }(),
	}))

	l := logger.RootLogger.Child(&logger.ChildLoggerOptions{Module: "main"})

	bot, err := telegram.CreateBot()
	if err != nil {
		l.FATAL("telegram bot didn't start: ", err)
	}

	bot.AddHandler("/ping", commands.Ping)
	bot.AddHandler("/info", commands.Info)

	bot.SetMyCommands()

	bot.Start()
}
