package main

import (
	"os"

	"github.com/dotenv-org/godotenvvault"
)

func main() {
	log("Starting...\n")
	err := godotenvvault.Load()
	if err != nil {
		logFatal("error loading dotenv", err)
	}

	logger := Logger(loggerArgs{prefix: "main"})
	logger.log("Hello World!\n")
	logger.f.log("We are in %s, token is %s", os.Getenv("ENV"), os.Getenv("TELEGRAM_TOKEN"))
}
