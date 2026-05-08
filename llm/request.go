package llm

import (
	"fmt"

	"github.com/sym01/htmlsanitizer"

	"github.com/d0ctr/bldbr-bot-go/llm/openai"
	"github.com/d0ctr/bldbr-bot-go/llm/types"
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

func SendRequest(model types.Model, messages []types.Message) (types.Message, error) {
	req := openai.BuildRequest(model, prompt, messages)

	if res, err := openai.SendRequest(model, req); err != nil {
		return types.Message{}, err
	} else if text, err := openai.FindFirstText(res); err != nil {
		return types.Message{}, err
	} else {
		text = sanitizeHtml(text)

		msg := types.FromText("", "", types.MESSAGE_ROLE_ASSISTANT, text)
		return msg, nil
	}
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
