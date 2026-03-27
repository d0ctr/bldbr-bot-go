package llm

import (
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func mapRole(r MessageRole) responses.EasyInputMessageRole {
	switch r {
	case MESSAGE_ROLE_USER:
		return responses.EasyInputMessageRoleUser
	case MESSAGE_ROLE_ASSISTANT:
		return responses.EasyInputMessageRoleAssistant
	case MESSAGE_ROLE_SYSTEM:
		return responses.EasyInputMessageRoleSystem
	}
	panic("unreachable")
}

func mapContentItem(c *MessageContent) responses.ResponseInputContentUnionParam{
	var inputText *responses.ResponseInputTextParam = nil
	var inputImage *responses.ResponseInputImageParam = nil
	switch c.t {
	case _MessageContentTypeMedia:
		inputImage = &responses.ResponseInputImageParam{
			ImageURL: openai.String(fmt.Sprintf("data:%s;base64,%s",c.media.mediaType, c.media.base64)),
		}
	case _MessageContentTypeText:
		inputText = &responses.ResponseInputTextParam {
			Text: c.text,
		}

	}
	return responses.ResponseInputContentUnionParam{
		OfInputText: inputText,
		OfInputImage: inputImage,
	}
}

func mapContentList(contents []*MessageContent) responses.EasyInputMessageContentUnionParam {
	if len(contents) == 1 && contents[0].t == _MessageContentTypeText {
		text := contents[0].text
		return responses.EasyInputMessageContentUnionParam{ OfString: openai.String(text) }
	}

	var list responses.ResponseInputMessageContentListParam

	for _, c := range contents {
		list = append(list, mapContentItem(c))
	}

	return responses.EasyInputMessageContentUnionParam{ OfInputItemContentList: list }
}

func mapMessages(messages []Message) responses.ResponseNewParamsInputUnion {
	var items []responses.ResponseInputItemUnionParam

	for _, message := range messages {
		item := responses.ResponseInputItemUnionParam{
			OfMessage: &responses.EasyInputMessageParam{
				Role: mapRole(message.role),
				Content: mapContentList(message.content),
			},
		}

		items = append(items, item)
	}
	return responses.ResponseNewParamsInputUnion{
		OfInputItemList: items,
	}
}
