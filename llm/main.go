package llm

import (
	"log/slog"
)

var logger = slog.With("component", "llm")

type MessageRole string

const (
	MESSAGE_ROLE_SYSTEM    MessageRole = "system"
	MESSAGE_ROLE_ASSISTANT MessageRole = "assistant"
	MESSAGE_ROLE_USER      MessageRole = "user"
)

type Message struct {
	author string
	role MessageRole
	content []*MessageContent
}

func NewMessage(author string, role MessageRole, content []*MessageContent) *Message {
	return &Message{ author, role, content }
}

func FromText(author string, role MessageRole, text string) *Message {
	content := []*MessageContent{ NewMessageContentText(text) }

	return NewMessage(author, role, content)
}

func (m Message) GetText() (string, bool) {
	for _, c := range m.content {
		if c.t == _MessageContentTypeText {
			return c.text, true
		}
	}

	return "", false
}

type MessageContentType string

const (
	_MessageContentTypeText  MessageContentType = "text"
	_MessageContentTypeMedia MessageContentType = "media"
)

type MessageContent struct {
	text string
	media MessageMedia
	t MessageContentType
}

func NewMessageContentText(text string) *MessageContent {
	return &MessageContent{ text, MessageMedia{}, _MessageContentTypeText }
}

func NewMessageContentMedia(fileId, base64, mediaType string) *MessageContent {
	return &MessageContent{ "", MessageMedia{ fileId, base64, mediaType }, _MessageContentTypeMedia }
}

type MessageMedia struct {
	fileId string
	base64 string
	mediaType string
}
