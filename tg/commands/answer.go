package commands

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/llm"
	"github.com/d0ctr/bldbr-bot-go/llm/types"
	"github.com/d0ctr/bldbr-bot-go/tg/utils"
)

var Answer = CommandDefinition{
	Name: "answer",
	Description: "{запрос?} сгенерировать текствой слоп",
	Handler: answer,
}

func answer(b *gotgbot.Bot, ctx *ext.Context) error {
	return llm.HandleTgCommand(b, ctx)
}

