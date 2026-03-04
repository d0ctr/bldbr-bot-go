package shared

import (
	"testing"

)

func TestGetValue(t *testing.T) {
	v, ok := GetValue(AHEGAO_API, "")

	if !ok {
		t.Error("expected `ok` to be true")
	}

	if len(v) == 0 {
		t.Errorf("expected a valid value, got [%s]", v)
	}
}
