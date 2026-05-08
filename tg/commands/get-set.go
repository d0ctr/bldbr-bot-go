package commands

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/shared"
	"github.com/d0ctr/bldbr-bot-go/tg/utils"
)

var Get = CommandDefinition{
	Name: "get",
	Description: "{название} вызвать запись из кэша",
	Handler: get,
}


var Set = CommandDefinition{
	Name: "set",
	Description: "{название} добавить запись в кэш",
	Handler: set,
}

var Lst = CommandDefinition{
	Name: "lst",
	Description: "список записей в кэше",
	Handler: lst,
}

func get(bot *gotgbot.Bot, ctx *ext.Context) error {
	name, ok := parseArgs(ctx.Message.Text, 1)[0]
	
	if !ok {
		sendErrorMsg(bot, ctx, "Не хватает названия записи.")
		return fmt.Errorf("get is missing an argument")
	}

	if err := validateName(bot, ctx, name); err != nil {
		return err
	}

	redis, rtx := shared.Redis()

	key := toKey(ctx, name)
	keyType, err := redis.Type(rtx, key).Result()
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при получении.", err)
		return fmt.Errorf("error getting value type from redis: %w", err)
	}

	var text string
	var media *Media

	switch keyType {
	case "none":
		sendErrorMsg(bot, ctx, "Такой записи пока не существует.")
		return fmt.Errorf("provided key does not exist in redis")
	case "hash":
		result, err := redis.HGetAll(rtx, key).Result()
		if err != nil {
			sendErrorMsg(bot, ctx, "Ошибка при получении.", err)
			return fmt.Errorf("error getting value from redis: %w", err)
		}

		var data _Data
		err = json.Unmarshal([]byte(result["data"]), &data)
		if err != nil {
			sendErrorMsg(bot, ctx, "Ошибка при получении.", err)
			return fmt.Errorf("error decoding value from redis: %w", err)
		}

		text = data.Text
		media = &Media{
			Id: data.GetMedia(),
			Type: data.Type,
		}
	default:
		text, err = redis.Get(rtx, key).Result()
		if err != nil {
			sendErrorMsg(bot, ctx, "Ошибка при получении.", err)
			return fmt.Errorf("error getting value from redis: %w", err)
		}
	}

	if media != nil {
		_, err = sendMedia(bot, ctx, text, media)
	} else {
		_, err = ctx.Message.Reply(bot, text, utils.GetDefaultMessageOpts())
	}
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при отправке.", err)
		return fmt.Errorf("error sending the message: %w", err)
	}


	return nil
}

func lst(bot *gotgbot.Bot, ctx *ext.Context) error {
	redis, rtx := shared.Redis()

	keys, err := redis.Keys(rtx, toKey(ctx, "*")).Result()
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при получении.", err)
		return fmt.Errorf("error getting keys from redis: %w", err)
	}

	if len(keys) == 0 {
		sendErrorMsg(bot, ctx, "В этом чате пока нет записей.")
		return fmt.Errorf("no keys with query [%s] found in redis", toKey(ctx, "*"))
	}

	for i, key := range keys {
		keys[i] = fmt.Sprintf(" - <code>%s</code>", strings.Split(key, ":")[2])
	}

	text :=  fmt.Sprintf("Записи доступные в этом чате:\n%s", strings.Join(keys, "\n"))
	_, err = ctx.Message.Reply(bot, text, utils.GetDefaultMessageOpts())
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при отправке.", err)
		return fmt.Errorf("error sending the message: %w", err)
	}

	return nil
}

func set(bot *gotgbot.Bot, ctx *ext.Context) error {
	redis, rtx := shared.Redis()

	message := ctx.Message.ReplyToMessage
	if message == nil {
		sendErrorMsg(bot, ctx, "Команда должна быть отправлена в ответ на другое сообщение.")
		return nil
	}

	name, ok := parseArgs(ctx.Message.Text, 1)[0]
	if !ok {
		sendErrorMsg(bot, ctx, "Не хватает название для записи, оно может состоять из букв, цифр и символов <code>_</code> и <code>-</code>.")
		return nil
	}

	if err := validateName(bot, ctx, name); err != nil {
		return err
	}

	text := message.OriginalHTML()
	if len(text) == 0 {
		text = message.OriginalCaptionHTML()
	}

	media := getMedia(message)

	data := _Data{
		Text: text,
	}

	if media.Type != "" {
		data.Type = media.Type
		if media.Type != MEDIA_TYPE_TEXT {
			data.Media = media.Id
		}
	}

	hash := make(map[string]any)
	if dataJson, err := json.Marshal(data); err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при сохранении", err)
		return err
	} else {
		hash["data"] = dataJson
	}

	if ctx.Message.From.Id != ctx.Message.Chat.Id {
		hash["owner"] = ctx.Message.From.Id
	}

	if _, err := redis.HSet(rtx, toKey(ctx, name), hash).Result(); err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при сохранении", err)
		return err
	}

	if _, err := ctx.Message.Reply(bot, fmt.Sprintf("Запись была сохранена под названием <code>%s</code>", name), utils.GetDefaultMessageOpts()); err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при отправке сообщения", err)
		return err
	}

	return nil
}

func validateName(bot *gotgbot.Bot, ctx *ext.Context, name string) error {
	valid := true
	for _, s := range name {
		if (!unicode.In(s, unicode.Letter, unicode.Number) && !slices.Contains([]rune{'-', '_'}, s)) {
			valid = false
			break
		}
	}

	if !valid {
		sendErrorMsg(bot, ctx, "Название гета может содержать только буквы, цифры, <code>-</code>, <code>_<code>.")
		return fmt.Errorf("get name contains unallowed character")
	}

	return nil

}

func toKey(ctx *ext.Context, name string) string {
	return fmt.Sprintf("%d:get:%s", ctx.Message.Chat.Id, name)
}

type _Data struct {
	Type        MediaType  `json:"type"`
	Owner       any         `json:"owner"`

	Text        string      `json:"text,omitempty"`

	// only one of the following
	Audio       string      `json:"audio,omitempty"`
	Animation   string      `json:"animation,omitempty"`
	Document    string      `json:"document,omitempty"`
	Photo       string      `json:"photo,omitempty"`
	Sticker     string      `json:"sticker,omitempty"`
	Video       string      `json:"video,omitempty"`
	VideoNote   string      `json:"video_note,omitempty"`
	Voice       string      `json:"voice,omitempty"`

	// a replacement for union-like data
	Media       string      `json:"media,omitempty"`
}

func (d *_Data) GetMedia() (string) {
	if d.Type == MEDIA_TYPE_TEXT || d.Type == "" {
		return ""
	} else if d.Media != "" {
		return d.Media
	}

	switch(d.Type) {
		case MEDIA_TYPE_AUDIO:       return d.Audio
		case MEDIA_TYPE_ANIMATION:   return d.Animation
		case MEDIA_TYPE_DOCUMENT:    return d.Document
		case MEDIA_TYPE_PHOTO:       return d.Photo
		case MEDIA_TYPE_STICKER:     return d.Sticker
		case MEDIA_TYPE_VIDEO:       return d.Video
		case MEDIA_TYPE_VIDEO_NOTE:  return d.VideoNote
		case MEDIA_TYPE_VOICE:       return d.Voice
	}

	return ""
}
