package commands

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"regexp"


	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/d0ctr/bldbr-bot-go/shared"
)

type CommandDefinition struct {
	Name, Description string
	Handler handlers.Response
}

var All = func() map[string]CommandDefinition {
	all_a := []CommandDefinition{Ping, Ahegao, Urban, Get, Set, Lst}

	all_m := make(map[string]CommandDefinition, len(all_a))
	
	for _, def := range all_a {
		all_m[def.Name] = def
	}

	return all_m
}()

var _WORDS_RE = regexp.MustCompile(` ([^ ]+)`)

func parseArgs(text string, limit uint) map[uint]string {
	matches := _WORDS_RE.FindAllStringSubmatch(text, -1)

	words := make(map[uint]string, len(matches))
	for i, match := range matches {
		words[uint(i)] = match[1]
	}

	if (limit > 0 && limit <= uint(len(words))) {
		args := make(map[uint]string, limit)

		i := uint(0);
		for ; i < limit - 1; i++ {
			args[i] = words[i]
		}

		lastArgSlice := make([]string, len(words) - int(limit) + 1)
		for ; i < uint(len(words)); i++ {
			lastArgSlice[i - limit + 1] = words[i]
		}
		args[limit - 1] = strings.Join(lastArgSlice, " ")

		return args
	}
	return words
}

func sendErrorMsg(bot *gotgbot.Bot, ctx *ext.Context, msg string, errs ...error) (*gotgbot.Message, error) {
	if len(errs) > 0 {
		msg = fmt.Sprintf("%s : \n<code>%s</code>", msg, errs[0].Error())
	}
	return ctx.Message.Reply(bot, msg, &shared.DEFAULT_MESSAGE_OPTS)
}

func handleHttpResponse(bot *gotgbot.Bot, ctx *ext.Context, entity string, r *http.Response, err error, statusCodes ...int) error {
	if len(statusCodes) == 0 {
		statusCodes = []int{http.StatusOK}
	}
	if err != nil || !slices.Contains(statusCodes, r.StatusCode) {
		if err == nil {
			err = fmt.Errorf("request to %s has failed with status [%s]", entity, r.Status)
		}
		sendErrorMsg(bot, ctx, "Ошибка при запросе", err)
		return fmt.Errorf("failed to get %s: %w", entity, err)
	}

	return nil
}

type _MediaType string;

const (
	MEDIA_TYPE_AUDIO       _MediaType = "audio"
	MEDIA_TYPE_ANIMATION   _MediaType = "animation"
	MEDIA_TYPE_DOCUMENT    _MediaType = "document"
	MEDIA_TYPE_PHOTO       _MediaType = "photo"
	MEDIA_TYPE_STICKER     _MediaType = "sticker"
	MEDIA_TYPE_VIDEO       _MediaType = "video"
	MEDIA_TYPE_VIDEO_NOTE  _MediaType = "video_note"
	MEDIA_TYPE_VOICE       _MediaType = "voice"

	// not an actual media type but a fallback type
	MEDIA_TYPE_TEXT        _MediaType = "text"
	
	// these could be supported but would require a more sophisticated reference than a simple file id
	// MEDIA_TYPE_MEDIA_GROUP _MediaType = "media_group"
	// MEDIA_TYPE_CONTACT     _MediaType = "contact"
	// MEDIA_TYPE_DICE        _MediaType = "dice"
	// MEDIA_TYPE_GAME        _MediaType = "game"
	// MEDIA_TYPE_INVOICE     _MediaType = "invoice"
	// MEDIA_TYPE_LOCATION    _MediaType = "location"
	// MEDIA_TYPE_POLL        _MediaType = "poll"
	// MEDIA_TYPE_VENUE       _MediaType = "venue"
)

type _Media struct {
	Id string
	Type _MediaType
}
func getMedia(message *gotgbot.Message) (media _Media) {
	if message.Audio != nil {
		media.Id   = message.Audio.FileId
		media.Type = MEDIA_TYPE_AUDIO
	} else if message.Animation != nil {
		media.Id   = message.Animation.FileId
		media.Type = MEDIA_TYPE_ANIMATION
	} else if message.Document != nil {
		media.Id   = message.Document.FileId
		media.Type = MEDIA_TYPE_DOCUMENT
	} else if message.Photo != nil {
		media.Id   = message.Photo[0].FileId
		media.Type = MEDIA_TYPE_PHOTO
	} else if message.Sticker != nil {
		media.Id   = message.Sticker.FileId
		media.Type = MEDIA_TYPE_STICKER
	} else if message.Video != nil {
		media.Id   = message.Video.FileId
		media.Type = MEDIA_TYPE_VIDEO
	} else if message.VideoNote != nil {
		media.Id   = message.VideoNote.FileId
		media.Type = MEDIA_TYPE_VIDEO_NOTE
	} else if message.Voice != nil {
		media.Id   = message.Voice.FileId
		media.Type = MEDIA_TYPE_VOICE
	}

	return media
}

func sendMedia(bot *gotgbot.Bot, ctx *ext.Context, text string, media *_Media) (*gotgbot.Message, error) {
	var msg *gotgbot.Message
	var err error
	switch media.Type {
	case MEDIA_TYPE_ANIMATION: 
		opts := gotgbot.SendAnimationOpts{
			Caption: text,
			ReplyParameters: &shared.DEFAULT_REPLY_PARAMETERS,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyAnimation(bot, input, &opts)
	case MEDIA_TYPE_AUDIO:
		opts := gotgbot.SendAudioOpts{
			Caption: text,
			ReplyParameters: &shared.DEFAULT_REPLY_PARAMETERS,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyAudio(bot, input, &opts)
	case MEDIA_TYPE_DOCUMENT:
		opts := gotgbot.SendDocumentOpts{
			Caption: text,
			ReplyParameters: &shared.DEFAULT_REPLY_PARAMETERS,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyDocument(bot, input, &opts)
	case MEDIA_TYPE_PHOTO:
		opts := gotgbot.SendPhotoOpts{
			Caption: text,
			ReplyParameters: &shared.DEFAULT_REPLY_PARAMETERS,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyPhoto(bot, input, &opts)
	case MEDIA_TYPE_STICKER:
		opts := gotgbot.SendStickerOpts{
			ReplyParameters: &shared.DEFAULT_REPLY_PARAMETERS,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplySticker(bot, input, &opts)
	case MEDIA_TYPE_VIDEO:
		opts := gotgbot.SendVideoOpts{
			Caption: text,
			ReplyParameters: &shared.DEFAULT_REPLY_PARAMETERS,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyVideo(bot, input, &opts)
	case MEDIA_TYPE_VIDEO_NOTE:
		opts := gotgbot.SendVideoNoteOpts{
			ReplyParameters: &shared.DEFAULT_REPLY_PARAMETERS,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyVideoNote(bot, input, &opts)
	case MEDIA_TYPE_VOICE:
		opts := gotgbot.SendVoiceOpts{
			Caption: text,
			ReplyParameters: &shared.DEFAULT_REPLY_PARAMETERS,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyVoice(bot, input, &opts)
	default:
		if text != "" {
			msg, err = ctx.Message.Reply(bot, text, &shared.DEFAULT_MESSAGE_OPTS)
		}

	}

	if msg == nil {
		err = fmt.Errorf("unsupported media type [%s] or/and no text to send", media.Type)
	}

	if err != nil {
		return nil, err
	}

	return msg, nil

}
