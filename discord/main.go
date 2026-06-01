package discord

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
	"d0ctr/bldbr-bot/shared"
)

var logger *slog.Logger
var session *discordgo.Session

type Error error

var NoSessionError Error = fmt.Errorf("no session")

func init() {
	logger = slog.With(shared.ComponentAttr("discord"))
	var err error

	if token, ok := shared.DISCORD_TOKEN.Get(); !ok {
		logger.Warn("'DISCORD_TOKEN' not found")
	} else if session, err = discordgo.New(fmt.Sprintf("Bot %s", token)); err != nil {
		logger.Error("failed to create bot", shared.ErrAttr(err))
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
	ImageUrl string
}

func GetEvents(guildId string) (result []Event, err error) {
	if session == nil {
		return result, NoSessionError
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

		if e.Image != "" {
			event.ImageUrl = fmt.Sprintf("https://cdn.discordapp.com/guild-events/%s/%s.png?size=512", e.ID, e.Image)
		}

		result = append(result, event)
	}

	slices.SortStableFunc(result, func(a, b Event) int {
		return a.Start.Compare(b.Start)
	})

	return result, err
}
