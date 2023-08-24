package commands

import "github.com/d0ctr/bldbr-bot-go/common/locale"

type infoCommand struct {
	*command
}

func (c infoCommand) Name() string {
	return c.name
}

func (c infoCommand) Telegram(ctx TelegramContext) error {
	senderId := ctx.Sender().ID
	chatId := ctx.Chat().ID
	chatType := ctx.Chat().Type
	l := c.telegramLog(ctx)

	result, err := locale.Get(locale.Language(ctx.Sender().LanguageCode), "command_info_answer", chatId, chatType, senderId)
	if err != nil {
		l.ERROR("error getting answer line: ", err)
	}

	return ctx.Reply(result)
}

var Info = infoCommand{
	command: &command{name: "info"},
}
