package commands

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/d0ctr/bldbr-bot-go/llm/types"
	"github.com/d0ctr/bldbr-bot-go/tg/utils"
)


var Model = CommandDefinition{
	Name: "model",
	Description: "{name?} изменить используемую модель",
	Handler: model,
}

func model(b *gotgbot.Bot, ctx *ext.Context) error {
	name, ok := parseArgs(ctx.EffectiveMessage.Text, 1)[0]

	available := &strings.Builder{}

	available.WriteString("Доступные модели:\n")
	for _, name := range types.GetModelNames() {
		fmt.Fprintf(available, "<code>  -</code> <code>%s</code>\n", name)
	}

	sendAvailable := func () {
		ctx.EffectiveMessage.Reply(b, available.String(), utils.GetDefaultMessageOpts())
	}

	if !ok {
		defer sendAvailable()
		cur, err := utils.GetChatValue(ctx, "llm-model")
		if err != nil {
			return err
		}

		if mdl, ok := types.GetOrDefault(cur); !ok {
			return fmt.Errorf("unknown model [%v]", cur)
		} else {
			fmt.Fprintf(available, "Текущая модель: <code>%s</code>", mdl.Name())
			return nil
		}
	} else {
		mdl, ok := types.GetOrDefault(name)
		if !ok {
			sendAvailable()
			return fmt.Errorf("unknown model [%s]", name)
		}

		if err := utils.SetChatValue(ctx, "llm-model", name); err != nil {
			return err
		}

		_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Модель изменена: <code>%s</code>", mdl.Name()), utils.GetDefaultMessageOpts())
		return err
	}
}
