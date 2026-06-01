package shared

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"os"
	"slices"
	"strings"
)

const (
	_COMPONENTS_DELIM = ":"
	_BAD_KEY          = "!BADKEY"
	_COMPONENT        = "component"
	_TEMPLATE         = "!TEMPLATE"
	_ERROR            = "error"
)

func init() {
	slog.SetDefault(NewLogger())
}

type _CustomHandler struct {
	inner *slog.JSONHandler
	attrs []slog.Attr
}

func TemplateAttr(value any) slog.Attr {
	return slog.Any(_TEMPLATE, value)
}

func ErrAttr(err error) slog.Attr {
	return slog.Any(_ERROR, err)
}

func ComponentAttr(name string) slog.Attr {
	return slog.String(_COMPONENT, name)
}

func newCustomHandler() *_CustomHandler {
	opts := &slog.HandlerOptions{ Level: slog.LevelDebug }
	innerHandler := slog.NewJSONHandler(os.Stdout, opts)
	return &_CustomHandler{ innerHandler, make([]slog.Attr, 0) }
}

func (h *_CustomHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.inner.Enabled(ctx, lvl)
}

func (h *_CustomHandler) WithAttrs(other []slog.Attr) slog.Handler {
	attrs := slices.Concat(h.attrs, other)
	return &_CustomHandler{ h.inner, attrs }
}

func (h *_CustomHandler) WithGroup(name string) slog.Handler {
	// return &customHandler{h.inner.WithGroup(name).(*slog.JSONHandler), h.attrs}
	return h
}

func (h *_CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	effectiveRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	var components []string
	var attrs []slog.Attr
	var args []string

	for attr := range h.Attrs(&r) {
		switch attr.Key {
		case _COMPONENT:
			components = append(components, attr.Value.String())
		case _TEMPLATE:
			args = append(args, fmt.Sprint(attr.Value.Any()))
		case _ERROR:
			attrs = append(attrs, slog.Any("error", attr.Value.Any()))
		default:
			attrs = append(attrs, attr)
		}
	}

	if len(components) > 0 {
		component := strings.Join(components, _COMPONENTS_DELIM)
		attrs = append(attrs, slog.String(_COMPONENT, component))
	}

	nextArg, _ := iter.Pull(slices.Values(args))

	splits := strings.Split(effectiveRecord.Message, "{}")
	b := strings.Builder{}
	b.WriteString(splits[0])
	for i := 1; i < len(splits); i += 1 {
		if arg, ok := nextArg(); ok {
			b.WriteString(arg)
		} else {
			b.WriteString("{}")
		}

		b.WriteString(splits[i])
	}
	effectiveRecord.Message = b.String()

	for arg, ok := nextArg(); ok; arg, ok = nextArg() {
		attrs = append(attrs, slog.String(_TEMPLATE, arg))
	}

	effectiveRecord.AddAttrs(attrs...)

	return h.inner.Handle(ctx, effectiveRecord)
}

func (h _CustomHandler) Attrs(r *slog.Record) iter.Seq[slog.Attr] {
	return func(yield func(slog.Attr) bool) {
		for _, attr := range h.attrs {
			if (!yield(attr)) {
				return
			}
		}

		r.Attrs(yield)
	}
}

func NewLogger() *slog.Logger {
	handler := newCustomHandler()
	logger := slog.New(handler).With(slog.String(_COMPONENT, "root"))
	return logger;
}
