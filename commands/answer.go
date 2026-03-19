package commands

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var Answer = CommandDefinition{
	Name: "answer",
	Description: "{запрос?} сгенерировать текствой слоп",
	Handler: answer,
}

func answer(bot *gotgbot.Bot, ctx *ext.Context) error {
	return nil
}

