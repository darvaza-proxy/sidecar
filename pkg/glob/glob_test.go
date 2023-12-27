package glob

import (
	"testing"
)

type captureCase struct {
	g   string
	sep []rune
	ok  bool
	mc  []captureMatchCase
}

type captureMatchCase struct {
	s  string
	ok bool
	c  []string
}

func tCC(s string, ok bool, c ...string) captureMatchCase {
	return captureMatchCase{
		s:  s,
		ok: ok,
		c:  c,
	}
}

func tC(g string, sep []rune, ok bool, mc ...captureMatchCase) captureCase {
	return captureCase{
		g:   g,
		sep: sep,
		ok:  ok,
		mc:  mc,
	}
}

func tD(g string, ok bool, m ...captureMatchCase) captureCase {
	return tC(g, []rune{'.'}, ok, m...)
}

func TestCapture(t *testing.T) {
	var cases = []captureCase{
		tD("*.local", true,
			tCC("me.local", true),
			tCC("local.you", false),
			tCC("me.local.you", false),
			tCC("me.you.local", false),
		),
		tD("(*).jpi.io", true,
			tCC("www.jpi.io", true, "www"),
		),
	}
	for _, tc := range cases {
		testCaptureCase(t, tc)
	}
}

func testCaptureCase(t *testing.T, tc captureCase) {
	g, err := Compile(tc.g, tc.sep...)
	switch {
	case tc.ok && err != nil:
		t.Errorf("%q: failed unexpectedly: %v", tc.g, err)
	case !tc.ok && err != nil:
		t.Logf("%q: failed as expected: %v", tc.g, err)
	case !tc.ok && err == nil:
		t.Errorf("%q: failed to fail", tc.g)
	default:
		testCaptureMatchCase(t, tc.g, g, tc.mc)
	}
}

func testCaptureMatchCase(t *testing.T, glob string, g *Glob, mc []captureMatchCase) {
	for _, mcc := range mc {
		fixture := mcc.s
		expected := mcc.c
		captured, ok := g.Capture(fixture)
		switch {
		case ok && !mcc.ok:
			t.Errorf("%q: %q: shouldn't have matched: %v",
				glob, fixture, captured)
		case !ok && mcc.ok:
			t.Errorf("%q: %q: should have matched: %v",
				glob, fixture, expected)
		case !ok && !mcc.ok:
			t.Logf("%q: %q: failed as expected", glob, fixture)
		case len(captured) != len(expected):
			t.Errorf("%q: %q: produced the wrong number of matches: %q instead of %q",
				glob, fixture, captured, expected)
		default:
			testCaptureMatchCaptures(t, glob, fixture, captured, expected)
		}
	}
}

func testCaptureMatchCaptures(t *testing.T, glob, fixture string, captured, expected []string) {
	for i := range captured {
		if captured[i] != expected[i] {
			t.Errorf("%q: %q: incorrect capture %v: %q != %q",
				glob, fixture, i+1, captured[i], expected[i])
		} else {
			t.Logf("%q: %q: capture %v: %q", glob, fixture, i+1, captured[i])
		}
	}
}
