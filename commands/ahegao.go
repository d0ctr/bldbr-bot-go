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
		ctx.EffectiveMessage.Reply(bot, "Эта команда не доступна", &shared.DEFAULT_MESSAGE_OPTS)
		return fmt.Errorf("command is enabled but its constraint is not satisfied (ahegao api url is unavailable)")
	}

	resp, err := http.Get(url)

	if err != nil {
		msg := fmt.Sprintf("Ошибка при получении картинки: %s", err.Error())
		ctx.EffectiveMessage.Reply(bot, msg, &shared.DEFAULT_MESSAGE_OPTS)
		return fmt.Errorf("failed to get ahegao url: %w", err)
	}

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Ошибка при получении картинки: %s", resp.Status)
		ctx.EffectiveMessage.Reply(bot, msg, &shared.DEFAULT_MESSAGE_OPTS)
		return fmt.Errorf("request to ahegao api has failed with status [%s]", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	data := make(map[string]string)
	if err := decoder.Decode(&data); err != nil {
		msg := fmt.Sprintf("Ошибка при получении картинки: %s", err.Error())
		ctx.EffectiveMessage.Reply(bot, msg, &shared.DEFAULT_MESSAGE_OPTS)
		return fmt.Errorf("failed to decode response from ahegao api: %w", err)
	}

	ahegao, ok := data["msg"]
	if !ok {
		return Ahegao(bot, ctx)
	}

	message, err := ctx.EffectiveMessage.ReplyPhoto(bot, gotgbot.InputFileByURL(ahegao), &shared.DEFAULT_PHOTO_OPTS)
	if err != nil {
		msg := fmt.Sprintf("Ошибка при отправке картинки: %s", err.Error())
		ctx.EffectiveMessage.Reply(bot, msg, &shared.DEFAULT_MESSAGE_OPTS)
		return fmt.Errorf("failed to send ahegao as url: %w", err)
	}

	logger.Debug("message sent", slog.Group("message", "id", message.MessageId, "chat_id", message.Chat.Id))

	return nil
}
