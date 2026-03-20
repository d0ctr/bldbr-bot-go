package llm

import (
	base64Enc "encoding/base64"
	"io"
	"net/http"
	"fmt"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func GetAuthorAndRole(bot *gotgbot.Bot, source *gotgbot.Message) (string, MessageRole) {
	var role MessageRole
	if source.From.Id == bot.Id {
		role = MESSAGE_ROLE_ASSISTANT
	} else {
		role = MESSAGE_ROLE_USER
	}

	author := fmt.Sprintf(`%s "%s" %s`, source.From.FirstName, source.From.Username, source.From.LastName)

	return author, role
}

func FromTgMessage(bot *gotgbot.Bot, source *gotgbot.Message) *Message {
	author, role := GetAuthorAndRole(bot, source)

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
