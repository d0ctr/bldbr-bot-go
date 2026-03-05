package main

import (
	"fmt"
	"context"
	"log/slog"
	"os"
	"slices"
	"strings"
)

const (
	_COMPONENTS_DELIM = ":"
	_BAD_KEY = "!BADKEY"
	COMPONENT = "component"
)

type _CustomHandler struct {
	inner *slog.JSONHandler
	attrs []slog.Attr
}

func newCustomHandler() *_CustomHandler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// for some reason time is already UTC, while slog's docs show local
		// ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		// 	if (a.Key ==  slog.TimeKey) {
		// 		return slog.Time(a.Key, a.Value.Time().UTC())
		// 	}
		// 	return a
		// },
	}
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
	consumeAttrs(h.attrs, &r, func(a slog.Attr) bool {
		switch a.Key {
		case COMPONENT:
			components = append(components, a.Value.String())
		case _BAD_KEY:
			var _args []any
			if v, ok := a.Value.Any().([]any); ok {
				_args = v
			} else {
				_args = []any{a.Value.Any()}
			}

			i := 0
			for arg := fmt.Sprint(_args[i]); i < len(_args) - 1; i++ {
				args = append(args, arg)
			}

			last := _args[i]
			if err, ok := last.(error); ok {
				attrs = append(attrs, slog.Any("error", err))
			} else {
				args = append(args, fmt.Sprint(last))
			}

		default:
			attrs = append(attrs, a)
		}

		return true
	})

	if (len(components) > 0) {
		component := strings.Join(components, _COMPONENTS_DELIM)
		effectiveRecord.Add(COMPONENT, component)
	}

	effectiveRecord.AddAttrs(attrs...)
	for _, arg := range args {
		effectiveRecord.Message = strings.Replace(effectiveRecord.Message, "{}", arg, 1)
	}

	return h.inner.Handle(ctx, effectiveRecord)
}

func consumeAttrs(attrs []slog.Attr, r *slog.Record, consumer func(slog.Attr) bool) {
	for _, attr := range attrs {
		if (!consumer(attr)) {
			return
		}
	}

	r.Attrs(consumer)
}

func NewLogger() *slog.Logger {
	handler := newCustomHandler()
	logger := slog.New(handler).With(slog.String(COMPONENT, "root"))
	return logger;
}
