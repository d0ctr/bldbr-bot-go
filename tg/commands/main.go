package commands

import (
	"fmt"
	"maps"
	"strings"
	"unicode"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"d0ctr/bldbr-bot/tg/utils"
)

type CommandDefinition struct {
	Name, Description string
	Response handlers.Response
}

var All = func() map[string]CommandDefinition {
	allA := []CommandDefinition{Ping, Ahegao, Urban, Get, Set, Lst, Voice, Answer, Model, Events, Info, Imagine, Context}

	allM := make(map[string]CommandDefinition, len(allA))
	
	for _, def := range allA {
		allM[def.Name] = def
	}

	return allM
}()

var sendErrorMsg = utils.SendErrorMsg
var withAction = utils.WithAction
var handleHttpResponse = utils.HandleHttpResponse

func parseArgs(text string, limit uint) map[uint]string {
	words := maps.Collect(func(yield func(uint, string) bool) {
		for i, field := range strings.Fields(text) {
			if i == 0 {
				continue
			} else {
				if !yield(uint(i - 1), field) {
					break
				}
			}
		}
	})

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

func trimCommand(s string) string {
	if !strings.HasPrefix(s, "/") {
		return s
	}

	space := strings.IndexFunc(s, unicode.IsSpace)

	if space == -1 {
		return ""
	}

	return s[space-1:]
}

func joinText(messages ...*gotgbot.Message) string {
	builder := strings.Builder{}

	for _, msg := range messages {
		if msg == nil {
			continue
		}

		text := msg.GetText()
		text = trimCommand(text)

		if text != "" {
			fmt.Fprintf(&builder, "%v\n", text)
		}
	}

	return builder.String()
}
