package commands

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"


	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/shared"
)

func parseArgs(text string, limit uint) map[uint]string {
	_, text, ok := strings.Cut(text, " ")
	args := make(map[uint]string)
	if !ok {
		return args
	}
	var i uint = 0
	for word := range strings.SplitSeq(text, " ") {
		if i < limit {
			if word == "" {
				continue
			}
			args[i] = word
			i += 1
		} else if last, ok := args[i]; ok {
			args[i] = strings.Join([]string{ last, word }, " ")
		} else {
			args[i] = word
		}
	}

	return args
}

func parseArgsDep(text string, limit uint) map[uint]string {
	var WORDS_RE = regexp.MustCompile(` ([^ ]+)`)

	matches := WORDS_RE.FindAllStringSubmatch(text, -1)

	words := make(map[uint]string, len(matches))
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
	}
	return words
}

func sendErrorMsg(bot *gotgbot.Bot, ctx *ext.Context, msg string, errs ...error) (*gotgbot.Message, error) {
	if len(errs) > 0 {
		msg = fmt.Sprintf("%s : \n<code>%s</code>", msg, errs[0].Error())
	}
	return ctx.Message.Reply(bot, msg, &shared.DEFAULT_MESSAGE_OPTS)
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

func withAction(bot *gotgbot.Bot, ctx *ext.Context, action string) func() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			ctx.Message.Chat.SendAction(bot, action, nil)
			<-ticker.C
		}

	}()

	return ticker.Stop
}
