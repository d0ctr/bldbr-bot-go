package openai

import (
	"fmt"
	"strings"

	"d0ctr/bldbr-bot/llm/types"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func mapRole(r types.MessageRole) responses.EasyInputMessageRole {
	switch r {
	case types.MESSAGE_ROLE_USER:
		return responses.EasyInputMessageRoleUser
	case types.MESSAGE_ROLE_ASSISTANT:
		return responses.EasyInputMessageRoleAssistant
	case types.MESSAGE_ROLE_SYSTEM:
		return responses.EasyInputMessageRoleSystem
	}

	panic("unreachable")
}

func mapContentItem(c types.MessageContent) responses.ResponseInputContentUnionParam{
	var inputText *responses.ResponseInputTextParam = nil
	var inputImage *responses.ResponseInputImageParam = nil

	c.OfText(func(text string) {
		inputText = &responses.ResponseInputTextParam {
			Text: text,
		}
	})

	c.OfMedia(func(mediaType, base64, fileId string) {
		inputImage = &responses.ResponseInputImageParam{
			ImageURL: openai.String(fmt.Sprintf("data:%s;base64,%s", mediaType, base64)),
		}
	})

	return responses.ResponseInputContentUnionParam{
		OfInputText: inputText,
		OfInputImage: inputImage,
	}
}

func mapContentList(contents []types.MessageContent, textOnly bool) (content responses.EasyInputMessageContentUnionParam) {
	if textOnly {
		var textBuilder = strings.Builder{}

		for _, c := range contents {
			c.OfText(func(text string) {
				fmt.Fprintf(&textBuilder, "%v\n", text)
			})
		}

		content.OfString = openai.String(textBuilder.String())

	} else {
		var list responses.ResponseInputMessageContentListParam

		for _, c := range contents {
			list = append(list, mapContentItem(c))
		}
		content.OfInputItemContentList = list

	}

	return content
}

func mapMessages(messages []types.Message) responses.ResponseNewParamsInputUnion {
	var items []responses.ResponseInputItemUnionParam

	for i, message := range messages {
		var textOnly bool

		if i == 0 || message.Role() == types.MESSAGE_ROLE_ASSISTANT {
			textOnly = true
		} else {
			textOnly = false
		}

		item := responses.ResponseInputItemUnionParam{
			OfMessage: &responses.EasyInputMessageParam{
				Role: mapRole(message.Role()),
				Content: mapContentList(message.Content(), textOnly),
			},
		}

		items = append(items, item)
	}
	return responses.ResponseNewParamsInputUnion{
		OfInputItemList: items,
	}
}
