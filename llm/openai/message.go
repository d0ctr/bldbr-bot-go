package openai

import (
	"fmt"

	"github.com/d0ctr/bldbr-bot-go/llm/types"
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

func mapContentList(contents []types.MessageContent) responses.EasyInputMessageContentUnionParam {
	var list responses.ResponseInputMessageContentListParam

	for _, c := range contents {
		list = append(list, mapContentItem(c))
	}

	return responses.EasyInputMessageContentUnionParam{ OfInputItemContentList: list }
}

func mapMessages(messages []types.Message) responses.ResponseNewParamsInputUnion {
	var items []responses.ResponseInputItemUnionParam

	for i, message := range messages {
		inputMessage := &responses.EasyInputMessageParam{
			Role: mapRole(message.Role()),
			Content: mapContentList(message.Content()),
		}

		if i == 0 {
			fixFirstMessage(inputMessage)
		}

		item := responses.ResponseInputItemUnionParam{
			OfMessage: inputMessage,
		}

		items = append(items, item)
	}
	return responses.ResponseNewParamsInputUnion{
		OfInputItemList: items,
	}
}

func fixFirstMessage(inputMessage *responses.EasyInputMessageParam) {
	if inputMessage.Role == responses.EasyInputMessageRoleUser {
		return
	}

	if len(inputMessage.Content.OfInputItemContentList) == 0 {
		return
	}

	for _, content := range inputMessage.Content.OfInputItemContentList {
		// non-text content may be emitted only by a user
		if content.OfInputText == nil {
			inputMessage.Role = responses.EasyInputMessageRoleUser
			break
		}
	}
}
