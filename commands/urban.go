package commands

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	net_url "net/url" 
	"encoding/json"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/d0ctr/bldbr-bot-go/shared"
)

func Urban(bot *gotgbot.Bot, ctx *ext.Context) error {
	logger := slog.Default().With("component", "urban")

	url, ok := shared.GetValue(shared.URBAN_API, "")
	if !ok {
		ctx.EffectiveMessage.Reply(bot, "Эта команда недоступна", &shared.DEFAULT_MESSAGE_OPTS)
		return fmt.Errorf("command is enabled but its constraint is not satisfied (urban api url is unavailable)")
	}

	q, ok := parseArgs(ctx.EffectiveMessage.GetText(), 1)[0]
	if !ok {
		logger.Debug("no arguments were found, falling back to random term")
		q = ""
	}

	r, err := getUrbanDefinitions(url, q)
	if err = handleHttpResponse(bot, ctx, "urban dictionary", r, err); err != nil {
		return err
	}

	decoder := json.NewDecoder(r.Body)
	data := _UrbanDefinitions{}
	if err := decoder.Decode(&data); err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при получении определения", err)
		return fmt.Errorf("failed to decode response from urban api: %w", err)
	}

	definition := findBestDefinition(data)

	message, err := ctx.EffectiveMessage.Reply(bot, definition.String(), &shared.DEFAULT_MESSAGE_OPTS)
	if err != nil {
		sendErrorMsg(bot, ctx, "Ошибка при отправке сообщения", err)
		return fmt.Errorf("failed to send definition: %w", err)
	}

	logger.Debug("message sent", slog.Group("message", "id", message.MessageId, "chat_id", message.Chat.Id))

	return nil
}

const (
	_URBAN_DEFINITION_URL = "https://www.urbandictionary.com/define.php"
	_EOT = '\000'
)

var (
	_ANNOTATION_RE = regexp.MustCompile(`\[[^\]]+\]`)
	_URBAN_DEFINITION_TEMPLATE = template.Must(template.New("URBAN_DEFINITION").Parse("<a href=\"{{.Link}}\">{{.Word}}</a>\n\n" +
		"{{.Definition}}\n\n" +
		"<blockquote>{{.Example}}</blockquote>\n\n" + 
		"{{.Up}} 👍|👎 {{.Down}}" + string(_EOT),
	))
)

type _AnnotatedString string

func (s _AnnotatedString) String() string {
	return _ANNOTATION_RE.ReplaceAllStringFunc(string(s), func(match string) string {
		term := strings.Trim(match, "[]")
		url, _ := net_url.Parse(_URBAN_DEFINITION_URL)
		qp := url.Query()
		qp.Add("term", term)
		url.RawQuery = qp.Encode()
		return fmt.Sprintf(`<a href="%s">%s</a>`, url, term)
	})
}

type _UrbanDefinition struct {
	Word       string           `json:"word"`
	Definition _AnnotatedString `json:"definition"`
	Example    _AnnotatedString `json:"example"`
	Up         int              `json:"thumbs_up"`
	Down       int              `json:"thumbs_down"`
	Link       string           `json:"permalink"`
}

func (d *_UrbanDefinition) String() string {
	pr, pw := io.Pipe()
	r := bufio.NewReader(pr)

	go _URBAN_DEFINITION_TEMPLATE.Execute(pw, d)

	str, _ := r.ReadString(_EOT)

	return str
}

type _UrbanDefinitions struct {
	List []_UrbanDefinition `json:"list"`
}

func getUrbanDefinitions(base_url string, term string) (*http.Response, error) {
	var endpoint string
	if term != "" {
		endpoint = "define"
	} else {
		endpoint = "random"
	}

	url, _ := net_url.Parse(base_url)
	url = url.JoinPath(net_url.PathEscape(endpoint))
	if term != "" {
		p := url.Query()
		p.Add("term", term)
		url.RawQuery = p.Encode()
	}
	return http.Get(url.String())
}

func findBestDefinition(definitions _UrbanDefinitions) _UrbanDefinition {
	lst := definitions.List
	slices.SortFunc(lst, func (a, b _UrbanDefinition) int {
		return cmp.Compare(b.Up + b.Down, a.Up + a.Down)
	})

	for _, definition := range lst {
		if definition.Up > definition.Down {
			return definition
		}
	}

	return lst[0]
}

