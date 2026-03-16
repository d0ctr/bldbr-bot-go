package shared

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var DEFAULT_REPLY_PARAMETERS = gotgbot.ReplyParameters{
	AllowSendingWithoutReply: true,
}

var DEFAULT_MESSAGE_OPTS = gotgbot.SendMessageOpts{
	ParseMode: "HTML",
	ReplyParameters: &DEFAULT_REPLY_PARAMETERS,
}

var DEFAULT_PHOTO_OPTS = gotgbot.SendPhotoOpts{
	ParseMode: "HTML",
	ReplyParameters: &DEFAULT_REPLY_PARAMETERS,
}

