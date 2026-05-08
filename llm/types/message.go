package types

import (
	"fmt"
)

type MessageRole string

const (
	MESSAGE_ROLE_SYSTEM    MessageRole = "system"
	MESSAGE_ROLE_ASSISTANT MessageRole = "assistant"
	MESSAGE_ROLE_USER      MessageRole = "user"
)

type Message struct {
	id string
	author string
	role MessageRole
	content []MessageContent
}

func FromText[T ~int | ~int64 | string](id T, author string, role MessageRole, text string) Message {
	idStr := fmt.Sprint(id)

	content := []MessageContent{ NewMessageContentText(text) }

	return Message{idStr, author, role, content}
}

func NewMessage[T ~int | ~int64 | string](id T, author string, role MessageRole, content []MessageContent) Message {
	idStr := fmt.Sprint(id)

	return Message{idStr, author, role, content}
}

func (m Message) Id() string {
	return m.id
}

func (m Message) Text() (string, bool) {
	for _, c := range m.content {
		if c.t == _MessageContentTypeText {
			return c.text, true
		}
	}

	return "", false
}

func (m Message) Role() MessageRole {
	return m.role
}

func (m *Message) ChangeRole(new MessageRole) {
	m.role = new
}

func (m Message) Content() []MessageContent {
	return m.content
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

func NewMessageContentText(text string) MessageContent {
	return MessageContent{ text, MessageMedia{}, _MessageContentTypeText }
}

func NewMessageContentMedia(fileId, base64, mediaType string) MessageContent {
	return MessageContent{ "", MessageMedia{ fileId, base64, mediaType }, _MessageContentTypeMedia }
}

func (c MessageContent) OfText(yield func(text string)) {
	if c.t == _MessageContentTypeText {
		yield(c.text)
	}
}

func (c MessageContent) OfMedia(yield func(mediaType string, base64 string, fileId string)) {
	if c.t == _MessageContentTypeMedia {
		yield(c.media.mediaType, c.media.base64, c.media.fileId)
	}
}

type MessageMedia struct {
	fileId string
	base64 string
	mediaType string
}
