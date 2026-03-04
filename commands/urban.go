package commands

import (
	"fmt"
	"log/slog"
	// "net/http"
	// "encoding/json"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/shared"
)

func Urban(bot *gotgbot.Bot, ctx *ext.Context) error {
	logger := slog.Default().With("component", "urban")

	_, ok := shared.GetValue(shared.URBAN_API, "")
	if !ok {
		ctx.EffectiveMessage.Reply(bot, "Эта команда недоступна", &shared.DEFAULT_MESSAGE_OPTS)
		return fmt.Errorf("command is enabled but its constraint is not satisfied (urban api url is unavailable)")
	}

	_, ok = parseArgs(ctx.EffectiveMessage.GetText(), 1)[0]
	if !ok {
		logger.Debug("no arguments were found, falling back to random term")
	}

	return nil
}
