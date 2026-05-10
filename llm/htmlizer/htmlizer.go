package htmlizer

import (
	"container/list"
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
	"github.com/sym01/htmlsanitizer"
)

type htmlizer struct {
	listIndent *list.List
	p *parser.Parser
	html *html.Renderer
	md *md.Renderer
	s *htmlsanitizer.HTMLSanitizer
	lineBreak bool
}

func new() htmlizer {
	izer := htmlizer{}

	izer.p = parser.NewWithExtensions(
		parser.FencedCode |
		parser.Strikethrough |
		parser.SpaceHeadings |
		parser.OrderedListStart)

	izer.html = html.NewRenderer(html.RendererOptions{
		Flags: html.SkipImages,
		RenderNodeHook: izer.renderNodeHook,
	})

	izer.md = md.NewRenderer()

	izer.listIndent = list.New()

	izer.s = htmlsanitizer.NewHTMLSanitizer()

	izer.s.AllowList.Tags = []*htmlsanitizer.Tag{
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

	return izer
}

type tagPair struct {
	ot, ct string
}

var headings []tagPair = []tagPair{
	{ "<b>",    "</b>" },
	{ "<b><i>", "</i></b>" },
	{ "<i>",    "</i>" },
}

type indent struct {
	v uint
}

func (h *htmlizer) countListIndent(w io.Writer, _ *ast.List, entering bool) {
	if entering {
		n := &indent{ v: 0 }
		h.listIndent.PushBack(n)
		if h.listIndent.Len() > 1 {
			w.Write([]byte("\n"))
		}
	} else {
		h.listIndent.Remove(h.listIndent.Back())
		if h.listIndent.Len() == 0 {
			w.Write([]byte("\n"))
		}
	}

}

func (h *htmlizer) customListItem(w io.Writer, i *ast.ListItem, entering bool) {

	if entering {

		prefix := strings.Repeat(" ", 4 * (h.listIndent.Len()))
		w.Write([]byte(prefix))

		if i.ListFlags & ast.ListTypeOrdered != 0 {
			n, _ := h.listIndent.Back().Value.(*indent)
			n.v += 1

			fmt.Fprintf(w, "%d. ", n.v)
		} else {
			w.Write([]byte("› "))
		}
	} else {
		w.Write([]byte("\n"))
	}
}

func (htmlizer) customEmphasis(w io.Writer, n ast.Node, entering bool) bool {
	var opener string
	if entering {
		opener = "<"
	} else {
		opener = "</"
	}

	switch n.(type) {
	case *ast.Strong: 
		w.Write([]byte(opener + "b>"))
	case *ast.Emph:
		w.Write([]byte(opener + "i>"))
	case *ast.Del:
		w.Write([]byte(opener + "s>"))
	default:
		return false
	}

	return true
}

func (htmlizer) customHeading(w io.Writer, h *ast.Heading, entering bool) {
	level := min(h.Level, len(headings))

	if entering {
		br := true
		if prev := ast.GetPrevNode(h); prev == nil {
			br = false
		} else {
			switch prev.(type) {
			case *ast.Heading, *ast.Paragraph, *ast.List:
				br = false
			}
		}

		if br {
			w.Write([]byte("\n"))
		}
		w.Write([]byte(headings[level - 1].ot))
	} else {
		w.Write([]byte(headings[level - 1].ct))
		w.Write([]byte("\n"))
		if _, isNextParagraph := ast.GetNextNode(h).(*ast.Paragraph); !isNextParagraph {
			w.Write([]byte("\n"))
		}
	}
}

func (htmlizer) customBlockQuote(w io.Writer, entering bool) {
	if entering {
		w.Write([]byte("<blockquote expandable>"))
	} else {
		w.Write([]byte("</blockquote>\n"))
	}
}

func (htmlizer) newLineOnExit(w io.Writer, _ ast.Node, entering bool) {
	if !entering {
		w.Write([]byte("\n"))
	}
}

func (h *htmlizer) renderNodeHook(w io.Writer, n ast.Node, entering bool) (ast.WalkStatus, bool) {

	switch tn := n.(type) {
	case *ast.Heading: 
		h.customHeading(w, tn, entering)
	case *ast.BlockQuote:
		h.customBlockQuote(w, entering)
	case *ast.Hardbreak, *ast.HorizontalRule:
		h.newLineOnExit(w, n, entering)
	case *ast.List:
		h.countListIndent(w, tn, entering)
	case *ast.ListItem:
		h.customListItem(w, tn, entering)
	default:
		return ast.GoToNext, h.customEmphasis(w, n, entering)
	}

	return ast.GoToNext, true
}

func htmlize(text string, unsanitized bool, debug io.Writer) (string, error) {
	htmlizer := new()
	doc := htmlizer.p.Parse([]byte(text))

	if debug != nil {
		ast.Print(debug, doc)
	}

	textB := markdown.Render(doc, htmlizer.html)

	if unsanitized {
		return string(textB), nil
	} else if sanitized, err := htmlizer.s.Sanitize(textB); err != nil {
		return string(textB), err
	} else {
		return string(sanitized), nil
	}
}

// returns unsanitized text in case of an error
func Htmlize(text string) (string, error) {
	return htmlize(text, false, nil)
}
