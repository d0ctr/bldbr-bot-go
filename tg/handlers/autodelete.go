package handlers

import (
	"d0ctr/bldbr-bot/tg/utils"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func AutoDelete() ext.Handler {
	return handlers.NewMessage(
		func(msg *gotgbot.Message) bool {
			return msg.PinnedMessage != nil
		},
		func(b *gotgbot.Bot, ctx *ext.Context) error {
			_, err := ctx.Message.PinnedMessage.Delete(b, nil)
			return utils.FmtNoSendError("%w", err)
		},
	)
}
