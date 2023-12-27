package glob

import (
	"testing"
)

func tP(parts ...templatePart) *Template {
	return &Template{
		parts: parts,
	}
}

func tL(s string) templatePart {
	return templatePart{
		literal: s,
	}
}

func tI(i int) templatePart {
	return templatePart{
		index: i,
	}
}

func TestTemplate(t *testing.T) {
	cases := []struct {
		t string
		p *Template
	}{
		{"hello",
			tP(tL("hello"))},
		{"hello$",
			tP(tL("hello$"))},
		{"hello$1",
			tP(tL("hello"), tI(1))},
		{"hello${1", nil},
		{"hello${1}",
			tP(tL("hello"), tI(1))},
		{"hello${1}world",
			tP(tL("hello"), tI(1), tL("world"))},
		{"hello${1w}orld", nil},
		{"hello${1world", nil},
		{"hello${world", nil},
		{"hello${0}world", nil},
		{"hello${-3}world", nil},
		{"hello${-3wo}rld", nil},
		{"hello${-3world", nil},
		{"hello${1}world${2}from${3}",
			tP(tL("hello"), tI(1), tL("world"), tI(2), tL("from"), tI(3))},
		{"a$b$c$5gh",
			tP(tL("a$b$c"), tI(5), tL("gh"))},
	}

	for _, tc := range cases {
		p, err := CompileTemplate(tc.t)
		switch {
		case err != nil && tc.p != nil:
			t.Errorf("ERROR: %q: failed unexpectedly: %v", tc.t, err)
		case err != nil && tc.p == nil:
			t.Logf("%q: failed successfully: %v", tc.t, err)
		case tc.p == nil:
			t.Errorf("ERROR: %q: failed to fail: %v", tc.t, p)
		case !tc.p.Equal(p):
			t.Errorf("ERROR: %q: produced the wrong result: %v != %v",
				tc.t, p, tc.p)
		default:
			t.Logf("%q: %v", tc.t, p)
		}
	}
}
