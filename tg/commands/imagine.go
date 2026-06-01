package commands

import (
	"context"
	"d0ctr/bldbr-bot/services"
	"d0ctr/bldbr-bot/tg/utils"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/openai/openai-go/v3"
)

var Imagine = CommandDefinition{
	Name: "imagine",
	Description: "{текст?} сгенерировать изображение",
	Response: imagine,
}

func imagine(b *gotgbot.Bot, ctx *ext.Context) error {
	const imgN = 1

	prompt := joinText(ctx.EffectiveMessage, ctx.EffectiveMessage.ReplyToMessage)

	if prompt == "" {
		utils.SendErrorMsg(b, ctx, "Добавь к команде запрос и/или отправь её в ответ на другое сообщение.")
		return utils.FmtNoSendError("command is missing a query")
	}

	client := services.OpenAi()

	action := utils.WithAction(b, ctx, gotgbot.ChatActionUploadPhoto)
	defer action()

	res, err := client.Images.Generate(
		context.Background(),
		openai.ImageGenerateParams{
			Prompt: prompt,
			N: openai.Int(imgN),
			Background: openai.ImageGenerateParamsBackgroundAuto,
			Model: "gpt-image-2",
			Moderation: openai.ImageGenerateParamsModerationLow,
			OutputFormat: openai.ImageGenerateParamsOutputFormatJPEG,
			Quality: openai.ImageGenerateParamsQualityAuto,
			// ResponseFormat: openai.ImageGenerateParamsResponseFormatB64JSON, // theo nly format for chatgpt models
		})

	if err != nil {
		return fmt.Errorf("error generating an image: %w", err)
	}

	inputPhotos := make([]gotgbot.InputMedia, imgN)
	for i, img := range res.Data {
		reader := strings.NewReader(img.B64JSON)
		decoder := base64.NewDecoder(base64.StdEncoding, reader)
		inputFile := gotgbot.InputFileByReader(fmt.Sprintf("%v.png", i), decoder)
		inputPhotos[i] = gotgbot.InputMediaPhoto{
			Media: inputFile,
		}
	}

	_, err = ctx.Message.ReplyMediaGroup(b, inputPhotos, &gotgbot.SendMediaGroupOpts{})

	return err
}
