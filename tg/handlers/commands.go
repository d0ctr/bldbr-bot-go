package handlers

import (
	"d0ctr/bldbr-bot/tg/commands"
	"iter"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func Commands() iter.Seq[ext.Handler] {
	return func(yield func(ext.Handler) bool) {
		for _, command := range commands.All {
			if !yield(handlers.NewCommand(command.Name, command.Response)) {
				break
			}
		}
	}
}
