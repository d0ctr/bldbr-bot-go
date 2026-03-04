package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"testing"
)

func mockLogger() (*slog.Logger, *bufio.Reader) {
	pr, pw := io.Pipe()
	r := bufio.NewReader(pr)
	testHandler := slog.NewJSONHandler(
		pw,
		&slog.HandlerOptions{ Level: slog.LevelDebug },
	)	
	h := &customHandler{ testHandler, make([]slog.Attr, 0)}
	return slog.New(h), r
}

func readJsonLine(r *bufio.Reader, t *testing.T) map[string]any {
	bytes, err := r.ReadBytes('\n')
	if err != nil {
		t.Errorf("failed to read log bytes: %s", err)
	}

	var data map[string]any
	if json.Unmarshal(bytes, &data) != nil {
		t.Errorf("failed to parse log line: %s", err)
	}
	return data
}

func TestMultipleComponents(t *testing.T) {
	l, r := mockLogger()

	l = l.With(COMPONENT, "1")

	// to avoid blocking, logging should be perfomed in a routine
	go l.Info(t.Name(), COMPONENT, "2")

	line:= readJsonLine(r, t)

	component := line[COMPONENT]
	if component != "1:2" {
		t.Errorf("expected [%s], actual [%s]", "1:2", component)
	}

}

func TestNoComponents(t *testing.T) {
	l, r := mockLogger()

	go l.Info(t.Name())

	line := readJsonLine(r, t)
	if component, ok := line[COMPONENT]; ok {
		t.Errorf("expected no components, actual [%s]", component)
	}
}

func TestTemplateArgs_asSlice(t *testing.T) {
	expected := "hello, 1!"
	l, r := mockLogger()

	go l.Info("{}, {}!", []any{"hello", 1})

	line := readJsonLine(r, t)

	msg := line["msg"]
	if msg != expected {
		t.Errorf("expected [%s], actual [%s]", expected, msg)
	}
}

func TestTemplateArgs_asSingleInt(t *testing.T) {
	expected := "hello, 1!"
	l, r := mockLogger()

	go l.Info("hello, {}!", 1)

	line := readJsonLine(r, t)

	msg := line["msg"]
	if msg != expected {
		t.Errorf("expected [%s], actual [%s]", expected, msg)
	}
}

func TestTemplateArgs_asSingleErr(t *testing.T) {
	l, r := mockLogger()

	go l.Info(t.Name(), fmt.Errorf("test"))

	line := readJsonLine(r, t)

	msg, err := line["msg"], line["error"]
	if msg != t.Name() {
		t.Errorf("expected [%s], actual [%s]", t.Name(), msg)
	}
	if err != "test" {
		t.Errorf("extected [%s], actual [%s]", "test", err)
	}
}

func TestTemplateArgs_asSliceWithError(t *testing.T) {
	expected := "hello, 1!"
	l, r := mockLogger()

	go l.Info("hello, {}!", []any{1, fmt.Errorf("test")})

	line := readJsonLine(r, t)

	msg, err := line["msg"], line["error"]
	if msg != expected {
		t.Errorf("expected [%s], actual [%s]", expected, msg)
	}
	if err != "test" {
		t.Errorf("extected [%s], actual [%s]", "test", err)
	}
}
