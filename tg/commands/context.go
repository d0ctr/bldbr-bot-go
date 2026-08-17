package commands

import (
	"d0ctr/bldbr-bot/llm"
	"d0ctr/bldbr-bot/tg/utils"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var Context = CommandDefinition{
	Name: "context",
	Description: "контекст сообщения или сообщение в виде контекстной ноды",
	Response: contextCommand,
}

func contextCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.Message.ReplyToMessage == nil {
		_, err := ctx.Message.Reply(b, "Работает только при ответе на другое сообщение.", utils.GetDefaultMessageOpts())
		return err
	}

	context := llm.StringifyContext(b, ctx.Message.ReplyToMessage)

	text := fmt.Sprintf("<pre><code class=\"language-json\">%v</code></pre>", context)
	_, err := ctx.Message.Reply(b, text, utils.GetDefaultMessageOpts())

	return err
}
