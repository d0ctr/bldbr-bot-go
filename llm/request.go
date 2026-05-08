package llm

import (
	"context"
	"fmt"
	"slices"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/sym01/htmlsanitizer"

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

func parseResponse(res *responses.Response, err error) (Message, error) {
	if err != nil {
		return Message{}, fmt.Errorf("api error: %v", err)
	}

	if len(res.Output) == 0 {
		return Message{}, fmt.Errorf("received an error: %v", res.Error)
	}

	for _, item := range res.Output {
		if item.Type == "message" && item.Status == "completed" {
			for _, content := range item.Content {
				if content.Type == "output_text" {
					text := sanitizeHtml(content.Text)
					result := FromText("", "", MESSAGE_ROLE_ASSISTANT, text)
					return result, nil
				}
			}
		}
	}

	return Message{}, fmt.Errorf("no valid response: %v", res.Output)
}

func sendRequestOpenAi(model Model, messages []Message) (Message, error) {
	fixedMessages := fixMessages(messages)

	client := shared.OpenAi()
	if client == nil {
		return Message{}, fmt.Errorf("openai service is not available")
	}

	body := responses.ResponseNewParams{
		Instructions: openai.String(prompt),
		Store: openai.Bool(false),
		Input: OpenAi.mapMessages(fixedMessages),
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
		Model: model.name,
	}

	return parseResponse(client.Responses.New(context.Background(), body))
}

func sendRequestGrok(model Model, messages []Message) (Message, error) {
	fixedMessages := fixMessages(messages)

	client := shared.XAi()
	if client == nil {
		return Message{}, fmt.Errorf("xai service is not available")
	}

	xSearchTool := responses.ToolParamOfCustom("x_search")
	xSearchTool.OfCustom.Type = "x_search"

	webSearchTool := responses.ToolParamOfCustom("web_search")
	webSearchTool.OfCustom.Type = "web_search"

	body := responses.ResponseNewParams{
		Instructions: openai.String(prompt),
		Store: openai.Bool(false),
		Input: OpenAi.mapMessages(fixedMessages),
		Tools: []responses.ToolUnionParam{ xSearchTool, webSearchTool },
		Reasoning: responses.ReasoningParam{
			Effort: responses.ReasoningEffortLow,
			Summary: openai.ReasoningSummaryAuto,
		},
		Model: model.name,
	}

	return parseResponse(client.Responses.New(context.Background(), body))
}

func SendRequest(model Model, messages []Message) (Message, error) {
	var result Message
	var err error

	switch model.provider {
		case MODEL_PROVIDER_XAI: result, err = sendRequestGrok(model, messages);
		case MODEL_PROVIDER_OPENAI: result, err = sendRequestOpenAi(model, messages);
		default: return Message{}, fmt.Errorf("undefined model provider [%v]", model.provider);
	}

	return result, err
}

func containsMediaContent(contents []MessageContent) bool {
	return slices.ContainsFunc(contents, func(content MessageContent) bool {
		return content.t == _MessageContentTypeMedia
	})
}

func fixMessages(messages []Message) []Message {
	// if len(messages) == 1 or first message contains media
	// -> then role must always be a user role

	var fixed []Message
	for _, message := range messages {
		if len(message.content) == 0 {
			continue
		}
		fixed = append(fixed, message)
	}

	if len(fixed) == 1 || containsMediaContent(fixed[0].content)  {
		fixed[0].role = MESSAGE_ROLE_USER
	}

	return fixed
}

var sanitizer *htmlsanitizer.HTMLSanitizer

func init() {
	sanitizer = htmlsanitizer.NewHTMLSanitizer()

	sanitizer.AllowList.Tags = []*htmlsanitizer.Tag{
		{ Name: "b",          Attr: []string{},               URLAttr: []string{} },
		{ Name: "strong",     Attr: []string{},               URLAttr: []string{} },
		{ Name: "i",          Attr: []string{},               URLAttr: []string{} },
		{ Name: "em",         Attr: []string{},               URLAttr: []string{} },
		{ Name: "u",          Attr: []string{},               URLAttr: []string{} },
		{ Name: "ins",        Attr: []string{},               URLAttr: []string{} },
		{ Name: "s",          Attr: []string{},               URLAttr: []string{} },
		{ Name: "strike",     Attr: []string{},               URLAttr: []string{} },
		{ Name: "del",        Attr: []string{},               URLAttr: []string{} },
		{ Name: "span",       Attr: []string{ "tg-spoiler" }, URLAttr: []string{} },
		{ Name: "tg-spoiler", Attr: []string{},               URLAttr: []string{} },
		{ Name: "a",          Attr: []string{},               URLAttr: []string{"href"} },
		{ Name: "code",       Attr: []string{ "class" },      URLAttr: []string{} },
		{ Name: "pre",        Attr: []string{},               URLAttr: []string{} },
		{ Name: "blockquote", Attr: []string{ "expandable" }, URLAttr: []string{} },
	}
}

func sanitizeHtml(text string) string {
	if fixed, err := sanitizer.SanitizeString(fmt.Sprintf("<p>%s</p>", text)); err != nil {
		return text
	} else {
		return fixed
	}
}
