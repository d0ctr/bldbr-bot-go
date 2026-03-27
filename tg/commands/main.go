package commands

import (
	"regexp"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/d0ctr/bldbr-bot-go/tg/utils"
)

type CommandDefinition struct {
	Name, Description string
	Handler handlers.Response
}

var All = func() map[string]CommandDefinition {
	all_a := []CommandDefinition{Ping, Ahegao, Urban, Get, Set, Lst, Voice, Answer}

	all_m := make(map[string]CommandDefinition, len(all_a))
	
	for _, def := range all_a {
		all_m[def.Name] = def
	}

	return all_m
}()

var sendErrorMsg = utils.SendErrorMsg
var withAction = utils.WithAction
var handleHttpResponse = utils.HandleHttpResponse

func parseArgs(text string, limit uint) map[uint]string {
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
