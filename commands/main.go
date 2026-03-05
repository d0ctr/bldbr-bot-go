package commands

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"regexp"


	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/shared"
)

var _WORDS_RE = regexp.MustCompile(` (\w+)`)

func parseArgs(text string, limit uint) map[uint]string {
	matches := _WORDS_RE.FindAllStringSubmatch(text, -1)

	words := make(map[uint]string)
	for i, match := range matches {
		words[uint(i)] = match[1]
	}

	if (limit > 0 && limit <= uint(len(words))) {
		args := make(map[uint]string, limit)

		i := uint(0);
		for ; i < limit - 1; i++ {
			args[i] = words[i]
		}

		lastArgSlice := make([]string, len(words) - int(limit) + 1)
		for ; i < uint(len(words)); i++ {
			lastArgSlice[i - limit + 1] = words[i]
		}
		args[limit - 1] = strings.Join(lastArgSlice, " ")

		return args
	} else {
		return words
	}
}

func sendErrorMsg(bot *gotgbot.Bot, ctx *ext.Context, msg string, errs ...error) (*gotgbot.Message, error) {
	if len(errs) > 0 {
		msg = fmt.Sprintf("%s : \n<code>%s</code>", msg, errs[0].Error())
	}
	return ctx.EffectiveMessage.Reply(bot, msg, &shared.DEFAULT_MESSAGE_OPTS)
}

func handleHttpResponse(bot *gotgbot.Bot, ctx *ext.Context, entity string, r *http.Response, err error, statusCodes ...int) error {
	if len(statusCodes) == 0 {
		statusCodes = []int{http.StatusOK}
	}
	if err != nil || !slices.Contains(statusCodes, r.StatusCode) {
		if err == nil {
			err = fmt.Errorf("request to %s has failed with status [%s]", entity, r.Status)
		}
		sendErrorMsg(bot, ctx, "Ошибка при запросе", err)
		return fmt.Errorf("failed to get %s: %w", entity, err)
	}

	return nil
}
