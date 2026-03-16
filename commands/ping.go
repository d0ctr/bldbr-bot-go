package commands

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/shared"
)

var Ping = CommandDefinition{
	Name: "ping",
	Description: "pong",
	Handler: ping,
}

func ping(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.Message.Reply(bot, "<code>pong</code>", &shared.DEFAULT_MESSAGE_OPTS)
	return err
}

