package llm

import (
	"log/slog"

	"github.com/d0ctr/bldbr-bot-go/llm/htmlizer"
	"github.com/d0ctr/bldbr-bot-go/llm/openai"
	"github.com/d0ctr/bldbr-bot-go/llm/types"
)

const prompt = `# DESCRIPTION

You are a helpful assistant contained in a Telegram bot. Your name is Bilderberg Butler.
Answer requests with expertise in the field required by the request.
Make your answers brief but explain more if requested by the user.
Always answer in the same language as the request.
`

func SendRequest(model types.Model, messages []types.Message) (types.Message, error) {
	req := openai.BuildRequest(model, prompt, messages)

	if res, err := openai.SendRequest(model, req); err != nil {
		return types.Message{}, err
	} else if text, err := openai.FindFirstText(res); err != nil {
		return types.Message{}, err
	} else {
		text, err = htmlizer.Htmlize(text)
		if err != nil {
			slog.Error("error sanitizing html", err)
		}

		msg := types.FromText("", "", types.MESSAGE_ROLE_ASSISTANT, text)
		return msg, nil
	}
}

