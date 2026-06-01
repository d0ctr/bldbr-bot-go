package services

import (
	"d0ctr/bldbr-bot/shared"
	"log/slog"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func XAi() *openai.Client {
	return xAiClient
}

var xAiClient *openai.Client 

func init() {
	token, ok := shared.XAI_TOKEN.Get()
	if !ok {
		slog.Error("'XAI_TOKEN' is not found")
		return
	}

	c := openai.NewClient(
		option.WithAPIKey(token),
		option.WithBaseURL("https://api.x.ai/v1"),
	)
	xAiClient =  &c
}
