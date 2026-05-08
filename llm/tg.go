package llm

import (
	base64Enc "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	tg "github.com/d0ctr/bldbr-bot-go/tg/utils"
)

func GetMessageParams(bot *gotgbot.Bot, source *gotgbot.Message) (id string, author string, role MessageRole) {
	id = strconv.FormatInt(source.MessageId, 10)

	if source.From.Id == bot.Id {
		role = MESSAGE_ROLE_ASSISTANT
	} else {
		role = MESSAGE_ROLE_USER
	}

	author = fmt.Sprintf(`%s "%s" %s`, source.From.FirstName, source.From.Username, source.From.LastName)

	return id, author, role
}

func FromTgMessage(bot *gotgbot.Bot, source *gotgbot.Message) Message {
	id, author, role := GetMessageParams(bot, source)

	var content []MessageContent

	if source.GetText() != "" {
		content = append(content, NewMessageContentText(source.GetText()))
	}

	if len(source.Photo) > 0 {
		image := slices.MinFunc(source.Photo, func(a, b gotgbot.PhotoSize) int {
			return int((a.Height + a.Width) - (b.Height + b.Width))
		})

		if file, err := bot.GetFile(image.FileId, nil); err != nil {
			logger.Error("failed to get file", err)
		} else if r, err := http.Get(file.URL(bot, nil)); err != nil {
			logger.Error("failed to download the file", err)
		} else if r.StatusCode != 200 {
			logger.Error("failed to download the file with status code [{}]", r.Status)
		} else if bytes, err := io.ReadAll(r.Body); err != nil {
			logger.Error("failed to read the file", err)
		} else {
			base64 := base64Enc.StdEncoding.EncodeToString(bytes)
			mediaType := r.Header.Get("Content-Type")
			if mediaType == "" || !strings.Contains(mediaType, "image"){
				mediaType = "image/jpeg"
			}

			messageMedia := NewMessageContentMedia(file.FileId, base64, mediaType)

			content = append(content, messageMedia)
		}

	}

	return Message{ id, author, role, content }
}

func (_Heap) GetTgTree(chatId int64) *Tree {
	treeId := strconv.FormatInt(chatId, 10)

	if tree, ok := Heap.tg[treeId]; ok {
		return tree
	} else {
		return nil
	}
}

func (_Heap) AddTgTree(chatId int64, tree *Tree) {
	treeId := strconv.FormatInt(chatId, 10)
	Heap.tg[treeId] = tree
}

func (t *Tree) GetTgNode(messageId int64) *TreeNode {
	nodeId := strconv.FormatInt(messageId, 10)

	if node, ok := t.nodes[nodeId]; ok {
		return node
	} else {
		return nil
	}
}

func (_Heap) GetTgTreeWithNode(bot *gotgbot.Bot, source *gotgbot.Message) (*Tree, *TreeNode) {
	chatId := source.Chat.Id
	tree := Heap.GetTgTree(chatId)
	var node *TreeNode

	if tree == nil {
		message := FromTgMessage(bot, source)
		node = NewTreeNode(message)
		tree = NewTree(node)
		
		Heap.AddTgTree(chatId, tree)
	} else {
		node = tree.GetTgNode(source.MessageId)
	}

	if node == nil {
		message := FromTgMessage(bot, source)
		node = NewTreeNode(message)
		tree.AddNode(node)
	}

	return tree, node
}

func HandleTgChain(bot *gotgbot.Bot, ctx *ext.Context) error {
	chatId := ctx.Message.Chat.Id
	var tree *Tree
	var node *TreeNode
	var model Model

	{
		modelName, _ := tg.GetChatValue(ctx, "llm-model")
		var ok bool
		model, ok = GetOrDefault(modelName)

		if !ok {
			return fmt.Errorf("no model named [%s]", modelName)
		}

	}

	if ctx.Message.ReplyToMessage != nil {
		tree, node = Heap.GetTgTreeWithNode(bot, ctx.Message.ReplyToMessage)
	}

	if tree == nil {
		tree = Heap.GetTgTree(chatId)
	}

	{
		prev := node
		message := FromTgMessage(bot, ctx.Message)
		node = NewTreeNode(message)

		if tree == nil {
			tree = NewTree(node)

			Heap.AddTgTree(ctx.Message.Chat.Id, tree)
		} else if prev == nil {
			tree.AddNode(node)
		} else {
			tree.AppendNode(prev, node)
		}
	}

	if node == nil {
		panic("this is a reply message it at least must have a new node created from the replyed message")
	}

	messages := tree.CollectMessages(node)

	endAction := tg.WithAction(bot, ctx, gotgbot.ChatActionTyping)
	defer endAction()

	r, err := SendRequest(model, messages)
	if err != nil {
		tg.SendErrorMsg(bot, ctx, "Ошибка при получении ответа", err)
		return fmt.Errorf("request resulted in an error: %w", err)
	}

	text, ok := r.GetText()
	if !ok {
		tg.SendErrorMsg(bot, ctx, "В ответе не было текста...")
		return fmt.Errorf("no text in the response: %v", r)
	}

	tgMessage, err := ctx.Message.ReplyMessage(bot, text, tg.GetDefaultMessageOpts())
	if err != nil {
		tg.SendErrorMsg(bot, ctx, "ошибка при отправке ответа", err)
	}

	prev := node
	message := FromTgMessage(bot, tgMessage)
	node = NewTreeNode(message)
	tree.AppendNode(prev, node)

	return nil
}
