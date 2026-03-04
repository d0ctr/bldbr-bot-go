package shared

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var DEFAULT_MESSAGE_OPTS = gotgbot.SendMessageOpts{
	ParseMode: "HTML",
};

var DEFAULT_PHOTO_OPTS = gotgbot.SendPhotoOpts{
	ParseMode: "HTML",
}

const (
	CONTEXT_LOGGER = "context_logger"
)
