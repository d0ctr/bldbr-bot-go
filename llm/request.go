package llm

import (
	"context"
	"fmt"
	"slices"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"

	"github.com/d0ctr/bldbr-bot-go/shared"
)

const prompt = `# DESCRIPTION

You are a helpful assistant contained in a Telegram bot. Your name is Bilderberg Butler.
Answer requests with expertise in the field required by the request.
Make your answers brief but explain more if requested by the user.
Always answer in the same language as the request.

# FORMATTING

Formatted text must always be formatted using HTML-like style. It includes following tags:

"""
<b>bold</b>, <strong>bold</strong>
<i>italic</i>, <em>italic</em>
<u>underline</u>, <ins>underline</ins>
<s>strikethrough</s>, <strike>strikethrough</strike>, <del>strikethrough</del>
<span class="tg-spoiler">spoiler</span>, <tg-spoiler>spoiler</tg-spoiler>
<b>bold <i>italic bold <s>italic bold strikethrough <span class="tg-spoiler">italic bold strikethrough spoiler</span></s> <u>underline italic bold</u></i> bold</b>
<a href="http://www.example.com/">inline URL</a>
<code>inline fixed-width code</code>
<pre>pre-formatted fixed-width code block</pre>
<pre><code class="language-python">pre-formatted fixed-width code block written in the Python programming language</code></pre>
<blockquote>Block quotation started\nBlock quotation continued\nThe last line of the block quotation</blockquote>
<blockquote expandable>Expandable block quotation started\nExpandable block quotation continued\nExpandable block quotation continued\nHidden by default part of the block quotation started\nExpandable block quotation continued\nThe last line of the block quotation</blockquote>
"""

- Only the tags mentioned above are currently supported.
- Regular text doesn't require a tag and can be placed around tags.
- New line is "\n".
- All "<", ">" and "&" symbols that are not a part of a tag or an HTML entity must be replaced with the corresponding HTML entities ("<" with "&lt;", ">" with "&gt;" and "&" with "&amp";).
- All numerical HTML entities are supported.
- The API currently supports only the following named HTML entities: "&lt;", "&gt;", "&amp;" and "&quot;".
- Use nested "pre" and "code" tags, to define programming language for "pre" entity.
- Programming language can't be specified for standalone "code" tags.
`

func SendRequest(messages []*Message) (*Message, error) {
	messages = fixMessages(messages)


	client := shared.OpenAi()
	if client == nil {
		return nil, fmt.Errorf("openai service is not available")
	}

	body := responses.ResponseNewParams{
		Instructions: openai.String(prompt),
		Store: openai.Bool(false),
		Input: mapMessages(messages),
		ToolChoice: responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsAuto),
		},
		Tools: []responses.ToolUnionParam{
			responses.ToolParamOfWebSearch(responses.WebSearchToolTypeWebSearch),
		},
		Reasoning: responses.ReasoningParam{
			Effort: responses.ReasoningEffortMedium,
			Summary: openai.ReasoningSummaryAuto,
		},
		Model: "gpt-5.4-nano",
	}

	ctx := context.Background()
	r, err := client.Responses.New(ctx, body)
	if err != nil {
		return nil, err
	}

	if len(r.Output) == 0 {
		return nil, fmt.Errorf("received an error: %v", r.Error)
	}

	for _, item := range r.Output {
		if item.Type == "message" && item.Status == "completed" {
			for _, content := range item.Content {
				if content.Type == "output_text" {
					result := FromText("", MESSAGE_ROLE_ASSISTANT, content.Text)
					return result, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no valid response: %v", r.Output)
}

func containsMediaContent(contents []*MessageContent) bool {
	return slices.ContainsFunc(contents, func(content *MessageContent) bool {
		return content.t == _MessageContentTypeMedia
	})
}

func fixMessages(messages []*Message) []*Message {
	// if len(messages) == 1 or first message contains media
	// -> then role must always be a user role

	if len(messages) == 1 || containsMediaContent(messages[0].content)  {
		messages[0].role = MESSAGE_ROLE_USER
	}

	return messages
}

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

func mapMessages(messages []*Message) responses.ResponseNewParamsInputUnion {
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
