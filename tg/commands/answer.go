package commands

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"d0ctr/bldbr-bot/llm"
)

var Answer = CommandDefinition{
	Name: "answer",
	Description: "{запрос?} сгенерировать текствой слоп",
	Handler: answer,
}

func answer(b *gotgbot.Bot, ctx *ext.Context) error {
	return llm.HandleTgCommand(b, ctx)
}

