package llm

import (
	base64Enc "encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"d0ctr/bldbr-bot/llm/types"
	tg "d0ctr/bldbr-bot/tg/utils"
)

func getMessageParams(b *gotgbot.Bot, source *gotgbot.Message) (id string, author string, role types.MessageRole) {
	id = strconv.FormatInt(source.MessageId, 10)

	if source.From.Id == b.Id {
		role = types.MESSAGE_ROLE_ASSISTANT
	} else {
		role = types.MESSAGE_ROLE_USER
	}

	author = fmt.Sprintf(`%s "%s" %s`, source.From.FirstName, source.From.Username, source.From.LastName)

	return id, author, role
}

func fromTgMessage(b *gotgbot.Bot, source *gotgbot.Message) (types.Message, bool) {
	logger := slog.With("component", "tg-to-llm")
	id, author, role := getMessageParams(b, source)

	var content []types.MessageContent

	{
		text := source.OriginalMDV2()
		if text == "" {
			text = source.OriginalCaptionMDV2()
		}

		cutFirstWord := false
		// cut command if present
		if strings.HasPrefix(text, "/") {
			cutFirstWord = true
		} else if strings.HasPrefix(text, "@") {
			var firstWord string
			strings.FieldsSeq(text)(func(word string) bool {
				firstWord = word
				return false
			})


			if firstWord[1:] == b.Username {
				cutFirstWord = true
			}
		}

		if cutFirstWord {
			start := strings.IndexFunc(text, unicode.IsSpace)
			// is not the last byte
			if start != -1 && start != (len(text) - 1) {
				text = text[start + 1:]
			} else {
				text = ""
			}
		}

		if source.Quote != nil {
			quoteBuilder := strings.Builder{}

			for line := range strings.Lines(source.Quote.Text) {
				fmt.Fprintf(&quoteBuilder, "> %s\n", line)
			}

			quoteBuilder.WriteString("\n")
			text = quoteBuilder.String() + text
		}

		if text != "" {
			content = append(content, types.NewMessageContentText(source.OriginalMDV2()))
		}
	}

	if len(source.Photo) > 0 {
		image := slices.MinFunc(source.Photo, func(a, b gotgbot.PhotoSize) int {
			return int((a.Height + a.Width) - (b.Height + b.Width))
		})

		if file, err := b.GetFile(image.FileId, nil); err != nil {
			logger.Error("failed to get file", err)
		} else if r, err := http.Get(file.URL(b, nil)); err != nil {
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

			messageMedia := types.NewMessageContentMedia(file.FileId, base64, mediaType)

			content = append(content, messageMedia)
		}

	}

	return types.NewMessage(id, author, role, content), len(content) > 0
}

func getTgTree(chatId int64) *Tree {
	treeId := strconv.FormatInt(chatId, 10)

	if tree, ok := heap.tg[treeId]; ok {
		return tree
	} else {
		return nil
	}
}

func addTgTree(chatId int64, tree *Tree) {
	treeId := strconv.FormatInt(chatId, 10)
	heap.tg[treeId] = tree
}

func (t *Tree) getTgNode(messageId int64) *TreeNode {
	nodeId := strconv.FormatInt(messageId, 10)

	if node, ok := t.nodes[nodeId]; ok {
		return node
	} else {
		return nil
	}
}

func getTgTreeWithNode(b *gotgbot.Bot, source *gotgbot.Message) (*Tree, *TreeNode) {
	chatId := source.Chat.Id
	tree := getTgTree(chatId)
	var node *TreeNode

	if tree == nil {
		if message, ok := fromTgMessage(b, source); ok {
			node = NewTreeNode(message)
			tree = NewTree(node)

			addTgTree(chatId, tree)
			return tree, node
		}
	} else {
		node = tree.getTgNode(source.MessageId)
	}

	if node == nil {
		if message, ok := fromTgMessage(b, source); ok {
			node = NewTreeNode(message)
			tree.AddNode(node)
		}
	}

	return tree, node
}

func ReplyResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	return RespondToTgMessage(false, b, ctx)
}

func CommandResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	return RespondToTgMessage(true, b, ctx)
}

func RespondToTgMessage(command bool, b *gotgbot.Bot, ctx *ext.Context) error{
	chatId := ctx.Message.Chat.Id
	var tree *Tree
	var node *TreeNode
	var model types.Model

	{
		modelName, _ := tg.GetChatValue(ctx, "llm-model")
		var ok bool
		model, ok = types.GetOrDefault(modelName)

		if !ok {
			return fmt.Errorf("no model named [%s]", modelName)
		}
	}

	if ctx.Message.ReplyToMessage != nil {
		tree, node = getTgTreeWithNode(b, ctx.Message.ReplyToMessage)
	}

	if tree == nil {
		tree = getTgTree(chatId)
	}

	if message, ok := fromTgMessage(b, ctx.Message); ok {
		prev := node
		node = NewTreeNode(message)
		if tree == nil {
			tree = NewTree(node)
			addTgTree(chatId, tree)
		} else if prev == nil {
			tree.AddNode(node)
		} else {
			tree.AppendNode(prev, node)
		}
	}

	if node == nil {
		if command {
			tg.SendErrorMsg(b, ctx, "Добавь к команде запрос и/или отправь её в ответ на другое сообщение.")
			return tg.FmtNoSendError("command is missing a query")
		} else {
			panic("this is a reply message, it at least must have a new node created from the replied message")
		}
	}

	messages := tree.CollectMessages(node)

	endAction := tg.WithAction(b, ctx, gotgbot.ChatActionTyping)
	defer endAction()

	var tgMessage *gotgbot.Message
	var err error

	if ctx.EffectiveChat.Type == "private" {
		tgMessage, err = streamProgress(model, messages, b, ctx)
	} else {
		tgMessage, err = sendOneOff(model, messages, b, ctx)
	}

	if err != nil {
		return err
	}

	if message, ok := fromTgMessage(b, tgMessage); ok {
		prev := node
		node = NewTreeNode(message)
		tree.AppendNode(prev, node)
	} else {
		return fmt.Errorf("failed to save context node, no tex in the message")
	}

	return nil
}

func sendOneOff(model types.Model, messages []types.Message, b *gotgbot.Bot, ctx *ext.Context) (message *gotgbot.Message, err error) {
	r, err := SendRequest(model, messages)
	if err != nil {
		tg.SendErrorMsg(b, ctx, "Ошибка при получении ответа", err)
		return nil, tg.FmtNoSendError("request resulted in an error: %w", err)
	}

	text, ok := r.Text()
	if !ok {
		tg.SendErrorMsg(b, ctx, "В ответе не было текста...")
		return nil, tg.FmtNoSendError("no text in the response: %v", r)
	}

	tgMessage, err := ctx.Message.ReplyMessage(b, text, tg.GetDefaultMessageOpts())
	if err != nil {
		return nil, fmt.Errorf("error on send: %w", err)
	}
	return tgMessage, nil
}

func streamProgress(model types.Model, messages []types.Message, b *gotgbot.Bot, ctx *ext.Context) (message *gotgbot.Message, err error) {
	chatId := ctx.EffectiveChat.Id
	draftId := rand.Int64()

	r, e := CreateStream(model, messages)
	var text string

	cooldown := time.NewTimer(0)
	defer cooldown.Stop()

	smoothener := time.NewTicker(200 * time.Millisecond)
	defer smoothener.Stop()
	for true {
		timer := time.NewTimer(10 * time.Second)
		defer timer.Stop()

		select {
		case err := <- e:
			if err != nil {
				tg.SendErrorMsg(b, ctx, "Ошибка при получении ответа", err)
				return nil, tg.FmtNoSendError("request resulted in an error: %w", err)
			}
		case msg, open := <- r:
			if !open {
				goto finish
			}

			var ok bool
			text, ok = msg.Text()

			if !ok {
				tg.SendErrorMsg(b, ctx, "В ответе не было текста...")
				return nil, tg.FmtNoSendError("no text in the response: %v", r)
			}
		case <-timer.C:
			slog.Debug("timed out waiting next message")
			continue
		}

		select {
		case <- cooldown.C:
			cooldown.Reset(0)
		default:
			continue
		}

		select {
		case <-smoothener.C:
			// noop
		default:
			continue
		}

		ok, err := b.SendMessageDraft(chatId, draftId, &gotgbot.SendMessageDraftOpts{
			Text: text,
			ParseMode: gotgbot.ParseModeHTML,
		})

		if err != nil && strings.Contains(err.Error(), "Too Many Requests") {
			cooldownDuration := time.Duration(parseCooldown(err)) * time.Second
			slog.Debug("going to fast, cooling down for {}s", cooldownDuration)

			cooldown.Reset(cooldownDuration * time.Second)
		} else if err != nil {
			slog.Warn("failed to send draft", slog.Any("chat", chatId), err)
		} else if !ok {
			slog.Warn("draft was not sent", slog.Any("chat", chatId))
		}
	}

finish:
	tgMessage, err := ctx.Message.ReplyMessage(b, text, tg.GetDefaultMessageOpts())
	if err != nil {
		return nil, fmt.Errorf("error on send: %w", err)
	}

	return tgMessage, nil
}

var cooldownRegex = regexp.MustCompile(`retry after (?<time>\d+)$`)

func parseCooldown(err error) int {
	text := err.Error()

	matches := cooldownRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		timeoutStr := matches[1]

		if v, err := strconv.Atoi(timeoutStr); err != nil {
			slog.Error("failed to parse cooldown from string {}", text)
		} else {
			return v
		}
	}

	return 0
}
