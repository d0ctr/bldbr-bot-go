package services

import (
	"d0ctr/bldbr-bot/shared"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func OpenAi() *openai.Client {
	if (openaiClient == nil) {
		openaiClient = getOpenaiClient()
	}
	return openaiClient
}

var openaiClient *openai.Client = nil 

func getOpenaiClient() *openai.Client {
	token, ok := shared.OPENAI_TOKEN.Get()
	if !ok {
		return nil
	}

	c := openai.NewClient(option.WithAPIKey(token))
	return &c
}
