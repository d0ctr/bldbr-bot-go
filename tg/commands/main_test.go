package commands

import (
	"fmt"
	"testing"
)

type parseArgsParameters struct {
	input string
	limit uint
	expected map[uint]string
}

func testParseArgs(t *testing.T, p parseArgsParameters) {
	args := parseArgs(p.input, p.limit)
	if len(p.expected) != len(args) {
		t.Errorf("expected [%v], actual [%v]", p.expected, args)
	}

	i := uint(0)
	for ; i < uint(len(p.expected)); i++ {
		expected, actual := p.expected[i], args[i]
		if (expected != actual) {
			t.Errorf("[%d] expected [%s], actual [%s]", i, expected, actual)
		}
	}
}

func TestParseArgs(t *testing.T) {
	tests := []parseArgsParameters {
		{ " bar baz bax", 0, map[uint]string{0: "bar", 1: "baz", 2: "bax"} },
		{ " bar baz bax", 1, map[uint]string{0: "bar baz bax"} },
		{ " bar baz bax", 2, map[uint]string{0: "bar", 1: "baz bax"} },
		{ " bar baz bax", 3, map[uint]string{0: "bar", 1: "baz", 2: "bax"} },
		{ "", 0, map[uint]string{} },
		{ "", 1, map[uint]string{} },
		{ "", 2, map[uint]string{} },
	}

	for _, parameters := range tests {
		testname := fmt.Sprintf("args=%d,limit=%d",len(parameters.expected), parameters.limit)
		t.Run(testname, func (t *testing.T) { testParseArgs(t, parameters) })
	}
}
