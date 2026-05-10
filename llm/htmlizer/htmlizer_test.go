package htmlizer

import (
	"fmt"
	"iter"
	"log"
	"strings"
	"testing"
)

func TestHtmlize(t *testing.T) {
	var i uint
	for input, expected := range testTuples {
		name := fmt.Sprintf("Test #%d", i + 1)
		t.Run(name, func(t *testing.T) { 
			actual, _ := htmlize(input, false, t.Output())
			j := 0
			expectedLines := strings.SplitSeq(expected, "\n")
			nextActualLine, stop := iter.Pull(strings.SplitSeq(actual, "\n"))
			defer stop()


			for expectedLine := range expectedLines {
				actualLine, ok := nextActualLine()
				if !ok {
					t.Errorf("%d: <EMPTY>", j + 1)
					t.Logf("Expected:%s", expectedLine)
					return
				}

				if strings.Compare(actualLine, expectedLine) != 0 {
					t.Errorf("%d: <FAIL>", j + 1)
					t.Logf("  Actual   : %s", actualLine)
					// t.Logf("  Actual(b): %v", []byte(actualLine))
					t.Logf("Expected   : %s", expectedLine)
					// t.Logf("Expected(b): %v", []byte(expectedLine))
					return
				}

				t.Logf("%d: %s", j + 1, actualLine)
				j += 1
			}

			if actualLine, ok := nextActualLine(); ok {
				t.Errorf("%d: <EXTRA>", j + 1)
				t.Logf("Actual:%s", actualLine)
			}
		})

		i += 1
		
	}
}

// Test resources {{{

var mdInputs []string = []string{
	"# Main Heading\n\n" +
	"## Secondary Heading\n\n" +
	"### Tertiary Heading\n\n" + // {{{}}}
	"Some paragraph with **bold text** and *italic text* and ***bold italic***.\n" +
	"Also _underscore italic_. And some ~~strikethrough~~ text or <ins>underlined</ins>.\n\n" +
	"Inline `code span` here.",

	"## Unordered\n\n" +
	"Asterisks tight:\n\n" +
	"* asterisk 1\n" +
	"* asterisk 2\n" +
	"* asterisk 3\n\n" +
	"## Ordered\n\n" +
	"Tight:\n\n" +
	"1. First\n" +
	"2. Second\n" +
	"3. Third\n\n" +
	"## Nested\n\n" +
	"1. First\n" +
	"2. Second:\n" +
	"    * Fee\n" +
	"    * Fie\n" +
	"    * Foe\n" +
	"3. Third",

	"Bla bla\n\n" +
	"``` oz\n" +
	"code blocks breakup paragraphs\n" +
	"```\n\n" +
	"Bla Bla\n\n" +
	"``` oz\n" +
	"multiple code blocks work okay\n" +
	"```\n\n" +
	"Bla Bla",
}

var htmlOutputs []string = []string{
	"<b>Main Heading</b>\n\n" +
	"<b><i>Secondary Heading</i></b>\n\n" +
	"<b><u>Tertiary Heading</u></b>\n" + 
	"Some paragraph with <b>bold text</b> and <i>italic text</i> and <b><i>bold italic</i></b>.\n" +
	"Also <i>underscore italic</i>. And some <s>strikethrough</s> text or <ins>underlined</ins>.\n\n" +
	"Inline <code>code span</code> here.\n",

	"<b><i>Unordered</i></b>\n" +
	"Asterisks tight:\n" +
	"    › asterisk 1\n" +
	"    › asterisk 2\n" +
	"    › asterisk 3\n\n" +
	"<b><i>Ordered</i></b>\n\n" +
	"Tight:\n" +
	"    1. First\n" +
	"    2. Second\n" +
	"    3. Third\n\n" +
	"<b><i>Nested</i></b>\n\n" +
	"    1. First\n" +
	"    2. Second:\n" +
	"        › Fee\n" +
	"        › Fie\n" +
	"        › Foe\n\n" +
	"    3. Third\n\n",

	"Bla bla\n\n" +
	"<pre><code class=\"language-oz\">" +
	"code blocks breakup paragraphs\n" +
	"</code></pre>\n\n" +
	"Bla Bla\n\n" +
	"<pre><code class=\"language-oz\">" +
	"multiple code blocks work okay\n" +
	"</code></pre>\n\n" +
	"Bla Bla\n",
}

var testTuples iter.Seq2[string, string] = func (yield func (string, string) bool) {
	if len(mdInputs) != len(htmlOutputs) {
		log.Panicf("len(mdInputs) != len(htmlOutputs)")
	}

	for i := range len(mdInputs) {
		if !yield(mdInputs[i], htmlOutputs[i]) {
			break
		}
	}
}
// }}}

// vim: foldmethod=marker
