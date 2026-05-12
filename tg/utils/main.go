package utils

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func GetDefaultReplyParameters() *gotgbot.ReplyParameters {
	return &gotgbot.ReplyParameters{
		AllowSendingWithoutReply: true,
	}
}

func GetDefaultMessageOpts() *gotgbot.SendMessageOpts {
	return &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyParameters: GetDefaultReplyParameters(),
	}
}

func GetDefaultPhotoOpts() *gotgbot.SendPhotoOpts {
	return &gotgbot.SendPhotoOpts{
		ParseMode: "HTML",
		ReplyParameters: GetDefaultReplyParameters(),
	}
}

func WithAction(bot *gotgbot.Bot, ctx *ext.Context, action string) func() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			ctx.Message.Chat.SendAction(bot, action, nil)
			<-ticker.C
		}

	}()

	return ticker.Stop
}

func SendErrorMsg(bot *gotgbot.Bot, ctx *ext.Context, msg string, errs ...error) (*gotgbot.Message, error) {
	if len(errs) > 0 {
		msg = fmt.Sprintf("%s : \n<code>%s</code>", msg, errs[0].Error())
	}
	return ctx.Message.Reply(bot, msg, GetDefaultMessageOpts())
}

func HandleHttpResponse(bot *gotgbot.Bot, ctx *ext.Context, entity string, r *http.Response, err error, statusCodes ...int) error {
	if len(statusCodes) == 0 {
		statusCodes = []int{http.StatusOK}
	}
	if err != nil || !slices.Contains(statusCodes, r.StatusCode) {
		if err == nil {
			err = fmt.Errorf("request to %s has failed with status [%s]", entity, r.Status)
		}
		SendErrorMsg(bot, ctx, "Ошибка при запросе", err)
		return fmt.Errorf("failed to get %s: %w", entity, err)
	}

	return nil
}

type NoSendError error

func FmtNoSendError(format string, args ...any) (err error) {
	return NoSendError(fmt.Errorf(format, args...))
}

