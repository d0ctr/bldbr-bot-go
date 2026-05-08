package shared

import (
	"log"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func XAi() *openai.Client {
	return xAiClient
}

var xAiClient *openai.Client 

func init() {
	token, ok := XAI_TOKEN.Get()
	if !ok {
		log.Printf("'XAI_TOKEN' is not found")
		return
	}
	log.Printf("'XAI_TOKEN' is found")

	c := openai.NewClient(
		option.WithAPIKey(token),
		option.WithBaseURL("https://api.x.ai/v1"),
	)
	xAiClient =  &c
}
