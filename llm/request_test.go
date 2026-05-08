package llm

import (
	"testing"
)

const htmlText = `
Конечно! Вот несколько <b>примеров</b> форматирования в HTML-подобном стиле для Telegram-бота:
<b>Жирный</b> и <i>курсив</i> и <u>подчёркнутый</u>.
<s>зачёркнутый</s>

<blockquote>Это цитата в одну строку.
Вторая строка цитаты.</blockquote>

<blockquote expanda
ble>Это <b>скрытая</b> часть по умолчанию.
Её можно раскрыть.</blockquote>

<span class="tg-spoiler">Это спойлер</span>

С <a href="http://www.example.com/">ссылкой</a> и <code>inline-code</code>.

<pre>Многострочный текст\nбез форматирования</pre>

<pre><code class="language-python">def hello():
    print("Hello, world!")</code></pre>
`

func testSanitizeHtml(t *testing.T, input string, expected string) {
	result := sanitizeHtml(input)

	if result != expected {
		t.Logf("Expected:\n%s", expected)
		t.Logf("Result:\n%s", result)

		t.Error("bad result")
	}
}

func TestSanitizeHtml(t *testing.T) {
	t.Run("sanitizeHtml", func(t *testing.T) { testSanitizeHtml(t, htmlText, htmlText) })
}
