package commands

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/llm"
	"github.com/d0ctr/bldbr-bot-go/shared"
)

var Answer = CommandDefinition{
	Name: "answer",
	Description: "{запрос?} сгенерировать текствой слоп",
	Handler: answer,
}

func answer(bot *gotgbot.Bot, ctx *ext.Context) error {
	var messages []*llm.Message

	if ctx.Message.ReplyToMessage != nil {
		messages = append(messages, llm.FromTgMessage(bot, ctx.Message.ReplyToMessage))
	}

	if query , ok := parseArgs(ctx.Message.GetText(), 1)[0]; ok {
		author, role := llm.GetAuthorAndRole(bot, ctx.Message)
		messages = append(messages, llm.FromText(author, role, query))
	}

	if len(messages) == 0 {
		sendErrorMsg(bot, ctx, "Добавь к команде запрос и/или отправь её в ответ на другое сообщение.")
		return fmt.Errorf("command is missing a query")
	}

	r, err := llm.SendRequest(messages)
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при получении ответа.")
		return fmt.Errorf("request resulted in an error: %w", err)
	}

	text, ok := r.GetText()
	if !ok {
		sendErrorMsg(bot, ctx, "В ответе не было текста...")
		return fmt.Errorf("no text in the response: %v", r)
	}

	_, err = ctx.Message.ReplyMessage(bot, text, &shared.DEFAULT_MESSAGE_OPTS)
	if err != nil {
		sendErrorMsg(bot, ctx, "ошибка при отправке ответа", err)
	}

	return nil
}

