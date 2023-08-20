package main

import (
	l "github.com/d0ctr/bldbr-bot-go/common/logger"

	"github.com/dotenv-org/godotenvvault"
)

var logger = l.CreateLogger(&l.LoggerOptions{Module: "main"})

func main() {
	logger.Log(l.LEVELINFO, "Starting...\n")
	err := godotenvvault.Load()
	if err != nil {
		logger.LogFatal("error loading dotenv", err)
	}

}
