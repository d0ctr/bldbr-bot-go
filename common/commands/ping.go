package commands

type pingCommand struct {
	*command
}

func (c pingCommand) Name() string {
	return c.name
}

func (c pingCommand) Telegram(ctx TelegramContext) error {
	c.telegramLog(ctx)

	return ctx.Reply("pong")
}

var Ping = pingCommand{
	command: &command{name: "ping"},
}
