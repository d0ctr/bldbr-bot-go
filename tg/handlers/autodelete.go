package handlers

import (
	"d0ctr/bldbr-bot/tg/utils"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func AutoDelete() ext.Handler {
	handler := handlers.NewMessage(
		func(msg *gotgbot.Message) bool {
			return msg.PinnedMessage != nil
		},
		func(b *gotgbot.Bot, ctx *ext.Context) error {
			if ok, err := ctx.Message.Delete(b, nil); err != nil {
				return utils.FmtNoSendError("%w", err)
			} else if !ok {
				return utils.FmtNoSendError("failed to delete a service message")
			} else {
				return nil
			}
		},
	)

	handler.AllowBot = true
	return handlers.NewNamedhandler("message_autodelete", handler)
}
