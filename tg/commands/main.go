package commands

import (
	"maps"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"d0ctr/bldbr-bot/tg/utils"
)

type CommandDefinition struct {
	Name, Description string
	Handler handlers.Response
}

var All = func() map[string]CommandDefinition {
	all_a := []CommandDefinition{Ping, Ahegao, Urban, Get, Set, Lst, Voice, Answer, Model, Events, Info}

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
