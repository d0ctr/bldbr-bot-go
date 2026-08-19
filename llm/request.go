package llm

import (
	"d0ctr/bldbr-bot/llm/openai"
	"d0ctr/bldbr-bot/llm/types"
)

const prompt = `# DESCRIPTION

You are a helpful assistant contained in a Telegram bot. Your name is Bilderberg Butler.
Answer requests with expertise in the field required by the request.
Make your answers brief but explain more if requested by the user.
Always answer in the same language as the request.
`

func SendRequest(model types.Model, messages []types.Message, conversationId string, cursor string) (types.Message, string, error) {
	// disable conversation-based routing and take any server
	if cursor == "" {
		conversationId = ""
	}

	req := openai.BuildRequest(model, prompt, messages, cursor)

	if res, err := openai.SendRequest(model, req, conversationId); err != nil {
		return types.Message{}, "", err
	} else if text, err := res.FindFirstText(); err != nil {
		return types.Message{}, "", err
	} else {
		msg := types.FromText("", "", types.MESSAGE_ROLE_ASSISTANT, text)
		return msg, cursor, nil
	}
}

