package commands

import (
	"fmt"
	"log/slog"
	"net/http"
	"encoding/json"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/shared"
)

func Ahegao(bot *gotgbot.Bot, ctx *ext.Context) error {
	logger := slog.Default().With("component", "ahegao")

	url, ok := shared.GetValue(shared.AHEGAO_API, "")
	if !ok {
		ctx.EffectiveMessage.Reply(bot, "Эта команда недоступна", &shared.DEFAULT_MESSAGE_OPTS)
		return fmt.Errorf("command is enabled but its constraint is not satisfied (ahegao api url is unavailable)")
	}

	r, err := http.Get(url)
	if err = handleHttpResponse(bot, ctx, "ahegao api", r, err); err != nil {
		return err
	}

	decoder := json.NewDecoder(r.Body)
	data := make(map[string]string)
	if err := decoder.Decode(&data); err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при получении картинки", err)
		return fmt.Errorf("failed to decode response from ahegao api: %w", err)
	}

	ahegao, ok := data["msg"]
	if !ok {
		return Ahegao(bot, ctx)
	}

	message, err := ctx.EffectiveMessage.ReplyPhoto(bot, gotgbot.InputFileByURL(ahegao), &shared.DEFAULT_PHOTO_OPTS)
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при отправлении картинки", err)
		return fmt.Errorf("failed to send ahegao as url: %w", err)
	}

	logger.Debug("message sent", slog.Group("message", "id", message.MessageId, "chat_id", message.Chat.Id))

	return nil
}
