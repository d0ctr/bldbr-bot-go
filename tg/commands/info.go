package commands

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"d0ctr/bldbr-bot/tg/utils"
)

var Info = CommandDefinition{
	Name: "info",
	Description: "информация о чате",
	Response: info,
}

func info(bot *gotgbot.Bot, ctx *ext.Context) error {
	b := &strings.Builder{}

	fmt.Fprintf(b, "id пользователя: <code>%v</code>\n", ctx.EffectiveSender.Id())
	if ctx.Message.Chat.Type != "private" {
		fmt.Fprintf(b, "id чата: <code>%v</code>\n", ctx.EffectiveChat.Id)
	}

	_, err := ctx.Message.Reply(bot, b.String(), utils.GetDefaultMessageOpts())
	return err
}

