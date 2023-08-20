package main

import (
	"github.com/d0ctr/bldbr-bot-go/common"

	"github.com/dotenv-org/godotenvvault"
)

var logger *common.Logger = common.CreateLogger("main", common.LEVELNOISE)

func main() {

	logger.Log(common.LEVELINFO, "Starting...\n")
	err := godotenvvault.Load()
	if err != nil {
		logger.LogFatal("error loading dotenv", err)
	}

	logger.Log(common.LEVELNOISE, "test")
	logger.Log(common.LEVELDEBUG, "test")
	logger.Log(common.LEVELINFO, "test")
	logger.Log(common.LEVELWARN, "test")
	logger.Log(common.LEVELERROR, "test")
	logger.Log(common.LEVELFATAL, "test")
}
