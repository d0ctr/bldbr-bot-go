package commands

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/llm"
	"github.com/d0ctr/bldbr-bot-go/tg/utils"
)

var Answer = CommandDefinition{
	Name: "answer",
	Description: "{запрос?} сгенерировать текствой слоп",
	Handler: answer,
}

func answer(bot *gotgbot.Bot, ctx *ext.Context) error {
	var tree *llm.Tree
	var node *llm.TreeNode

	if ctx.Message.ReplyToMessage != nil {
		tree, node = llm.Heap.GetTgTreeWithNode(bot, ctx.Message.ReplyToMessage)
	}

	if tree == nil {
		tree = llm.Heap.GetTgTree(ctx.Message.Chat.Id)
	}

	if query, ok := parseArgs(ctx.Message.GetText(), 1)[0]; ok {
		id, author, role := llm.GetMessageParams(bot, ctx.Message)
		prev := node
		message := llm.FromText(id, author, role, query)
		node = llm.NewTreeNode(message)

		if tree == nil {
			tree = llm.NewTree(node)

			llm.Heap.AddTgTree(ctx.Message.Chat.Id, tree)
		} else if prev == nil {
			tree.AddNode(node)
		} else {
			tree.AppendNode(prev, node)
		}
	}

	if node == nil {
		sendErrorMsg(bot, ctx, "Добавь к команде запрос и/или отправь её в ответ на другое сообщение.")
		return fmt.Errorf("command is missing a query")
	}

	messages := tree.CollectMessages(node)

	endAction := withAction(bot, ctx, gotgbot.ChatActionTyping)
	defer endAction()
	r, err := llm.SendRequest(messages)
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при получении ответа", err)
		return fmt.Errorf("request resulted in an error: %w", err)
	}

	text, ok := r.GetText()
	if !ok {
		sendErrorMsg(bot, ctx, "В ответе не было текста...")
		return fmt.Errorf("no text in the response: %v", r)
	}

	tgMessage, err := ctx.Message.ReplyMessage(bot, text, utils.GetDefaultMessageOpts())
	if err != nil {
		sendErrorMsg(bot, ctx, "ошибка при отправке ответа", err)
	}

	prev := node
	message := llm.FromTgMessage(bot, tgMessage)
	node = llm.NewTreeNode(message)
	tree.AppendNode(prev, node)

	return nil
}

