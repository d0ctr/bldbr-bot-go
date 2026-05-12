package discord

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/d0ctr/bldbr-bot-go/shared"
)

var logger *slog.Logger
var session *discordgo.Session

var NO_SESSION = fmt.Errorf("no discord session")

func init() {
	logger = slog.With("component", "discord")
	var err error

	if token, ok := shared.DISCORD_TOKEN.Get(); !ok {
		logger.Error("'DISCORD_TOKEN' not found")
	} else if session, err = discordgo.New(fmt.Sprintf("Bot %s", token)); err != nil {
		logger.Error("failed to create bot", err)
	}
}

func IsAvailable() bool {
	return session != nil
}

type Entity struct {
	Id string
	Name string
	Url string
}

type Event struct {
	Entity

	Guild Entity

	// not present for external events
	Channel *Entity 

	Description string
	IsActive bool
	Location string
	Start time.Time
	End *time.Time
}

func GetEvents(guildId string) (result []Event, err error) {
	if session == nil {
		return result, NO_SESSION
	}

	guild, err := session.Guild(guildId)
	if err != nil {
		return result, fmt.Errorf("failed to get relavant guild: %w", err)
	}

	channels := make(map[string]Entity)
	events, err := session.GuildScheduledEvents(guildId, false)

	if err != nil {
		return result, fmt.Errorf("failed to acquire events: %w", err)
	}

	for _, e := range events {
		var channel *Entity
		if e.ChannelID != "" {
			channelId := e.ChannelID
			channel, ok := channels[channelId]
			if !ok {
				if c, err := session.Channel(channelId); err != nil {
					channel = Entity{
						Id: c.ID,
						Name: c.Name,
					}

					channels[channelId] = channel
				}
			}
		}

		logger.Debug("got event [{}]", e.Name)
		event := Event{
			Entity: Entity{
				Id: e.ID,
				Name: e.Name,
			},
			Description: e.Description,
			IsActive: e.Status == discordgo.GuildScheduledEventStatusActive,
			Location: e.EntityMetadata.Location,
			Start: e.ScheduledStartTime,
			End: e.ScheduledEndTime,

			Guild: Entity{
				Id: guild.ID,
				Name: guild.Name,
			},

			Channel: channel,
		}

		result = append(result, event)
	}

	logger.Debug("returning {} events", len(result))
	return result, err
}
