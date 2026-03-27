package commands

import (
	"testing"
)

func TestUrbanDefinitionTemplate(t *testing.T) {
	def := _UrbanDefinition{
		Word: "word",
		Definition: "definition with [annotation]",
		Example: "example with [annotations]",
		Up: 1,
		Down: 0,
		Link: "https://example.org",
	}

	str := def.String()

	if str == "" {
		t.Error("expected a non-empty string")
	} else {
		t.Logf("got definition: %s", str)
	}

}
