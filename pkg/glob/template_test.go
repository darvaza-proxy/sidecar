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

type testReplaceCase struct {
	t  string
	r  string
	d  []string
	ok bool
}

func tR(template string, result string, data ...string) testReplaceCase {
	return testReplaceCase{
		t:  template,
		r:  result,
		d:  data,
		ok: true,
	}
}

func tRE(template string, data ...string) testReplaceCase {
	return testReplaceCase{
		t:  template,
		d:  data,
		ok: false,
	}
}

func TestReplace(t *testing.T) {
	var cases = []testReplaceCase{
		tR("foobar", "foobar"),
		tR("$1Foe", "oneFoe", "one"),
		tR("${1}Foe", "oneFoe", "one"),
		tR("${1}Foe", "Foe", "", "two"),
		tRE("${2}Foe", "one"),
		tR("Hello, ${2}", "Hello, world", "one", "world"),
	}

	for _, tc := range cases {
		p, err := CompileTemplate(tc.t)
		switch {
		case err != nil:
			t.Errorf("ERROR: %q: failed to compile: %v", tc.t, err)
		default:
			testReplaceResult(t, p, tc)
		}
	}
}

func testReplaceResult(t *testing.T, p *Template, tc testReplaceCase) {
	r, err := p.Replace(tc.d)
	switch {
	case !tc.ok && err == nil:
		t.Errorf("ERROR: %q+%q: failed to fail: %q",
			tc.t, tc.d, r)
	case !tc.ok && err != nil:
		t.Logf("%q+%q: failed successfully: %v",
			tc.t, tc.d, err)
	case tc.ok && err != nil:
		t.Errorf("ERROR: %q+%q: failed to replace: %v",
			tc.t, tc.d, err)
	case r != tc.r:
		t.Errorf("ERROR: %q+%q: produced the wrong result: %q != %q",
			tc.t, tc.d, r, tc.r)
	default:
		t.Logf("%q+%q => %q", tc.t, tc.d, r)
	}
}
