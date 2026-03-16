package shared

import (
	"testing"

)

func TestConfigKeyGet(t *testing.T) {
	v, ok := AHEGAO_API.Get()

	if !ok {
		t.Error("expected `ok` to be true")
	}

	if len(v) == 0 {
		t.Errorf("expected a valid value, got [%s]", v)
	}
}
