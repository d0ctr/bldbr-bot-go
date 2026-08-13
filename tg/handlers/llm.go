package handlers

import (
	"d0ctr/bldbr-bot/llm"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func ReplyToBot() ext.Handler {
	handler := handlers.NewMessage(
		func(msg *gotgbot.Message) bool {
			// is a reply
			return msg.ReplyToMessage != nil &&
			// a reply to a bot's message
			msg.ReplyToMessage.From.IsBot &&
			// not a reply from a bot
			!msg.From.IsBot &&
			// not escaped with '! '
			!strings.HasPrefix(msg.GetText(), "! ")
		},
		llm.ReplyResponse,
	)

	return handlers.NewNamedhandler("message_llm_reply", handler)
}
