package llm

import (
	base64Enc "encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"fmt"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var logger = slog.With("component", "llm")

type MessageRole string

const (
	MESSAGE_ROLE_SYSTEM    MessageRole = "system"
	MESSAGE_ROLE_ASSISTANT MessageRole = "assistant"
	MESSAGE_ROLE_USER      MessageRole = "user"
)

type (
	Message struct {
		author string
		role MessageRole
		content []*MessageContent
	}

	MessageContentType string

	MessageContent struct {
		text string
		media MessageMedia
		t MessageContentType
	}

	MessageMedia struct {
		fileId string
		base64 string
		mediaType string
	}

)

const (
	_MessageContentTypeText  MessageContentType = "text"
	_MessageContentTypeMedia MessageContentType = "media"
)

func NewMessageContentText(text string) *MessageContent {
	return &MessageContent{ text, MessageMedia{}, _MessageContentTypeText }
}

func NewMessageContentMedia(fileId, base64, mediaType string) *MessageContent {
	return &MessageContent{ "", MessageMedia{ fileId, base64, mediaType }, _MessageContentTypeMedia }
}

func NewMessage(author string, role MessageRole, content []*MessageContent) *Message {
	return &Message{ author, role, content }
}

func FromTgMessage(bot *gotgbot.Bot, source *gotgbot.Message) *Message {
	var role MessageRole
	if source.From.Id == bot.Id {
		role = MESSAGE_ROLE_ASSISTANT
	} else {
		role = MESSAGE_ROLE_USER
	}

	author := fmt.Sprintf(`%s "%s" %s`, source.From.FirstName, source.From.Username, source.From.LastName)

	var content []*MessageContent

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
			if mediaType == "" {
				mediaType = "image/jpeg"
			}

			messageMedia := NewMessageContentMedia(file.FileId, base64, mediaType)

			content = append(content, messageMedia)
		}

	}

	return NewMessage(author, role, content)
}
