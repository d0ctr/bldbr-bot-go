package commands

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/tg/utils"
)

type MediaType string;

const (
	MEDIA_TYPE_AUDIO       MediaType = "audio"
	MEDIA_TYPE_ANIMATION   MediaType = "animation"
	MEDIA_TYPE_DOCUMENT    MediaType = "document"
	MEDIA_TYPE_PHOTO       MediaType = "photo"
	MEDIA_TYPE_STICKER     MediaType = "sticker"
	MEDIA_TYPE_VIDEO       MediaType = "video"
	MEDIA_TYPE_VIDEO_NOTE  MediaType = "video_note"
	MEDIA_TYPE_VOICE       MediaType = "voice"

	// not an actual media type but a fallback type
	MEDIA_TYPE_TEXT        MediaType = "text"
	
	// these could be supported but would require a more sophisticated reference than a simple file id
	// MEDIA_TYPE_MEDIA_GROUP MediaType = "media_group"
	// MEDIA_TYPE_CONTACT     MediaType = "contact"
	// MEDIA_TYPE_DICE        MediaType = "dice"
	// MEDIA_TYPE_GAME        MediaType = "game"
	// MEDIA_TYPE_INVOICE     MediaType = "invoice"
	// MEDIA_TYPE_LOCATION    MediaType = "location"
	// MEDIA_TYPE_POLL        MediaType = "poll"
	// MEDIA_TYPE_VENUE       MediaType = "venue"
)

type Media struct {
	Id string
	Type MediaType
}

func getMedia(message *gotgbot.Message) (media Media) {
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

func sendMedia(bot *gotgbot.Bot, ctx *ext.Context, text string, media *Media) (*gotgbot.Message, error) {
	var msg *gotgbot.Message
	var err error
	replyParameters := utils.GetDefaultReplyParameters()
	switch media.Type {
	case MEDIA_TYPE_ANIMATION: 
		opts := gotgbot.SendAnimationOpts{
			Caption: text,
			ReplyParameters: replyParameters,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyAnimation(bot, input, &opts)
	case MEDIA_TYPE_AUDIO:
		opts := gotgbot.SendAudioOpts{
			Caption: text,
			ReplyParameters: replyParameters,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyAudio(bot, input, &opts)
	case MEDIA_TYPE_DOCUMENT:
		opts := gotgbot.SendDocumentOpts{
			Caption: text,
			ReplyParameters: replyParameters,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyDocument(bot, input, &opts)
	case MEDIA_TYPE_PHOTO:
		opts := gotgbot.SendPhotoOpts{
			Caption: text,
			ReplyParameters: replyParameters,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyPhoto(bot, input, &opts)
	case MEDIA_TYPE_STICKER:
		opts := gotgbot.SendStickerOpts{
			ReplyParameters: replyParameters,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplySticker(bot, input, &opts)
	case MEDIA_TYPE_VIDEO:
		opts := gotgbot.SendVideoOpts{
			Caption: text,
			ReplyParameters: replyParameters,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyVideo(bot, input, &opts)
	case MEDIA_TYPE_VIDEO_NOTE:
		opts := gotgbot.SendVideoNoteOpts{
			ReplyParameters: replyParameters,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyVideoNote(bot, input, &opts)
	case MEDIA_TYPE_VOICE:
		opts := gotgbot.SendVoiceOpts{
			Caption: text,
			ReplyParameters: replyParameters,
			ParseMode: gotgbot.ParseModeHTML,
		}

		input := gotgbot.InputFileByID(media.Id)

		msg, err = ctx.Message.ReplyVoice(bot, input, &opts)
	default:
		if text != "" {
			msg, err = ctx.Message.Reply(bot, text, utils.GetDefaultMessageOpts())
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
