package main

import (
	"os"
	"strings"

	"github.com/d0ctr/bldbr-bot-go/common"

	"github.com/dotenv-org/godotenvvault"
)

var logger *common.Logger = common.CreateLogger("main", common.LEVELDEBUG)

func main() {

	logger.Log(common.LEVELINFO, "Starting...\n")
	err := godotenvvault.Load()
	if err != nil {
		logger.LogFatal("error loading dotenv", err)
	}

}
