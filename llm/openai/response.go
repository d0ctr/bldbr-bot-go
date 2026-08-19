package openai

import (
	"fmt"

	"github.com/openai/openai-go/v3/responses"
)

type Response struct {
	oResponse *responses.Response
}

func (res Response) FindFirstText() (string, error) {
	if len(res.oResponse.Output) == 0 {
		return "", fmt.Errorf("received an error: %v", res.oResponse.Error)
	}

	for _, item := range res.oResponse.Output {
		if item.Type == "message" && item.Status == "completed" {
			for _, content := range item.Content {
				if content.Type == "output_text" {
					return content.Text, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no valid response: %v", res.oResponse.Output)
}

func (res Response) GetId() (string) {
	return res.oResponse.ID
}
