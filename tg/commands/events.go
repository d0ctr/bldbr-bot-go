package commands

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"d0ctr/bldbr-bot/discord"
	"d0ctr/bldbr-bot/shared"
	"d0ctr/bldbr-bot/tg/utils"
)

var Events = CommandDefinition{
	Name: "events",
	Description: "список будущих событий на дискорд сервере",
	Response: events.events,
}

type _events struct {}

var events = _events{}

func (_events) events(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	if !discord.IsAvailable() {
		sendErrorMsg(b, ctx, "Нет связи с дискордом")
		return nil
	}

	redis, rtx := shared.Redis()

	key := events.toKey(ctx)

	if v, err := redis.Exists(rtx, key).Result(); err != nil {
		return fmt.Errorf("error checking stored value: %w", err)
	} else if v == 0 {
		sendErrorMsg(b, ctx, "Этот чат ещё не подписан на эвенты, чтобы подписаться добавьте бота на сервер в дискорде и воспользуйтесь командой /subevents.")
		return nil
	}

	eventsMap := make(map[string][]discord.Event)

	if guildIds, err := redis.SMembers(rtx, key).Result(); err != nil {
		return fmt.Errorf("error getting stored value: %w", err)
	} else if len(guildIds) == 0 {
		sendErrorMsg(b, ctx, "Этот чат ещё не подписан на эвенты, чтобы подписаться добавьте бота на сервер в дискорде и воспользуйтесь командой /subevents.")
		return nil
	} else {
		for _, guildId := range guildIds {
			var events []discord.Event
			events, err = discord.GetEvents(guildId)

			if len(events) > 0 {
				eventsMap[guildId] = events
			}
		}

		if len(eventsMap) == 0 && err != nil {
			 sendErrorMsg(b, ctx, "Ошибка при запросе", err)
			 return nil
		} else if err != nil {
			slog.Error("silently logging request error", err)
		}
	}

	if len(eventsMap) == 0 {
		sendErrorMsg(b, ctx, "Нет запланированных событий")
		return nil
	}

	replyingTo := ctx.Message
	responseBuilder := &strings.Builder{}
	for _, events := range eventsMap {
		for i, event := range events {
			if i == 0 {
				fmt.Fprintf(responseBuilder, "<b>%s</b>\n", event.Guild.Name)
			}

			if i > 0 && i % 5 == 0 {
				if msg, err := replyingTo.Reply(b, responseBuilder.String(), utils.GetDefaultMessageOpts()); err != nil {
					replyingTo = msg
				} else {
					slog.Error("failed to send response", err)
				}

				responseBuilder.Reset()
			} else {
				fmt.Fprintf(responseBuilder, "<code>› </code><b>%s</b>\n", event.Name)
				fmt.Fprintf(
					responseBuilder,
					"<code>  </code>Начало: <tg-time unix=\"%d\" format=\"Dt\">%s (UTC)</tg-time>\n",
					event.Start.Unix(),
					event.Start.Format("January 2, 2006 at 15:04"),
				)

				if event.End != nil {
					duration := event.End.Sub(event.Start)
					var minutes int = int(duration.Minutes())
					var hours int
					hours, minutes = minutes / 60, minutes % 60

					fmt.Fprintf(
						responseBuilder,
						"<code>  </code>Продолжительность: %02d:%02d\n",
						hours,
						minutes,
					)
				}

				fmt.Fprint(responseBuilder, "\n")
			}
		}

		if text := responseBuilder.String(); text != "" {
			if msg, err := replyingTo.Reply(b, responseBuilder.String(), utils.GetDefaultMessageOpts()); err != nil {
				slog.Error("failed to send response", err)
			} else {
				replyingTo = msg
			}
		}

		responseBuilder.Reset()
	}

	return nil
}

func (_events) toKey(ctx *ext.Context) string {
	return fmt.Sprintf("telegram:%v:event_subscriber:guild_ids", ctx.Message.Chat.Id)
}
