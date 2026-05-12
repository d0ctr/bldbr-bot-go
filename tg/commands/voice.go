package commands

import (
	"context"
	"net/http"
	"fmt"
	"io"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/openai/openai-go/v3"
	"github.com/google/uuid"

	"d0ctr/bldbr-bot/shared"
)

var Voice = CommandDefinition{
	Name: "voice",
	Description: "Генерирует голосове сообщение из текста или аудио",
	Handler: voice,
}

func voice(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message.ReplyToMessage

	if msg == nil || (msg.Audio == nil && msg.Text == "")  {
		sendErrorMsg(bot, ctx, "Команду нужено отправит в ответ на текстовое или аудио сообщение.")
		return fmt.Errorf("command must be sent as a reply")
	}

	endAction := withAction(bot, ctx, gotgbot.ChatActionRecordVoice)
	defer endAction()

	var reader io.Reader
	if msg.Audio != nil {
		if _reader, err := fromFile(bot, msg.Audio.FileId); err != nil {
			sendErrorMsg(bot, ctx, "Не удалось получить аудио", err)
			return fmt.Errorf("error generating audio: %w", err)
		} else {
			reader = _reader
		}
	} else {
		client := shared.OpenAi()
		
		if client == nil {
			sendErrorMsg(bot, ctx, "Озвучка отдыхает.")
			return fmt.Errorf("openai is not available")
		} else if _reader, err := toSpeech(client, msg.Text); err != nil {
			sendErrorMsg(bot, ctx, "Не удалось сгенерировать аудио", err)
			return fmt.Errorf("error generating audio: %w", err)
		} else {
			reader = _reader
		}

	}
	name := uuid.NewString()
	input := gotgbot.InputFileByReader(name, reader)

	if _, err := msg.ReplyVoice(bot, input, nil); err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при отправке сообщения", err)
		return err
	}

	return nil
}

const prompt = `Delivery: Exaggerated and theatrical, with dramatic pauses, sudden outbursts, and gleeful cackling.

Voice: High-energy, eccentric, and slightly unhinged, with a manic enthusiasm that rises and falls unpredictably.

Tone: Excited, chaotic, and grandiose, as if reveling in the brilliance of a mad experiment.

Pronunciation: Sharp and expressive, with elongated vowels, sudden inflections, and an emphasis on big words to sound more diabolical.`

func toSpeech(client *openai.Client, text string) (io.Reader, error) {
	ctx := context.Background()
	body := openai.AudioSpeechNewParams{
		Model: openai.SpeechModelGPT4oMiniTTS,
		Input: text,
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatOpus,
		Voice: openai.AudioSpeechNewParamsVoiceUnion{
			OfString: openai.String(string(openai.AudioSpeechNewParamsVoiceStringMarin)),
		},
		Instructions: openai.String(prompt),
	}
	r, err := client.Audio.Speech.New(ctx, body)
	if err != nil {
		return nil, err
	} else if r.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response from tts api [%v]: %v", r.StatusCode, r.Body)
	}

	return r.Body, nil
}

func fromFile(bot *gotgbot.Bot, fileId string) (io.Reader, error) {
	file, err := bot.GetFile(fileId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load file info: %w", err)
	}

	url := file.URL(bot, nil)

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	} else if r.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response from telegram api [%v]: %v", r.StatusCode, r.Body)
	}

	return r.Body, nil
}
