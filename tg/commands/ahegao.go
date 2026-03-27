package commands

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/shared"
	"github.com/d0ctr/bldbr-bot-go/tg/utils"
)

var Ahegao = CommandDefinition{
	Name: "ahegao",
	Description: "рандомное ахегао",
	Handler: ahegao,
}

func ahegao(bot *gotgbot.Bot, ctx *ext.Context) error {
	logger := slog.Default().With("component", "ahegao")

	url, ok := shared.AHEGAO_API.Get()
	if !ok {
		ctx.Message.Reply(bot, "Эта команда недоступна", utils.GetDefaultMessageOpts())
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

	ahegao_url, ok := data["msg"]
	if !ok {
		return ahegao(bot, ctx)
	}

	message, err := ctx.Message.ReplyPhoto(bot, gotgbot.InputFileByURL(ahegao_url), utils.GetDefaultPhotoOpts())
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при отправлении картинки", err)
		return fmt.Errorf("failed to send ahegao as url: %w", err)
	}

	logger.Debug("message sent", slog.Group("message", "id", message.MessageId, "chat_id", message.Chat.Id))

	return nil
}
