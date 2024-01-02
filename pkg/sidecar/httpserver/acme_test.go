package httpserver

import "testing"

type testAcmeHTTP01PatternCase struct {
	s  string // URL.Path
	ok bool   // path match
	m  string // token match
}

func TestAcmeHTTP01Pattern(t *testing.T) {
	var cases = []testAcmeHTTP01PatternCase{
		{"", false, ""},
		{"/", false, ""},
		{"/foo", false, ""},
		{"/.well-known", false, ""},
		{"/.well-known/", false, ""},
		{"/.well-known/acme-cha", false, ""},
		{"/.well-known/acme-cha/", false, ""},
		{"/.well-known/acme-challenge", true, ""},
		{"/.well-known/acme-challenge/", true, ""},
		{"/.well-known/acme-challenge/foo", true, "foo"},
		{"/.well-known/acme-challenge/foo/bar", true, ""},
		{"../.well-known/acme-challenge/foo", false, ""},
	}

	for _, tc := range cases {
		testOneAcmeHTTP01Pattern(t, tc)
	}
}

func testOneAcmeHTTP01Pattern(t *testing.T, tc testAcmeHTTP01PatternCase) {
	m, ok := AcmeHTTP01Pattern.Capture(tc.s)
	switch {
	case ok && !tc.ok:
		t.Errorf("ERROR: %q: failed to fail", tc.s)
	case !ok && tc.ok:
		t.Errorf("ERROR: %q: failed unexpectedly", tc.s)
	case !ok && !tc.ok:
		t.Logf("%q: failed successfully", tc.s)
	case len(m) != 1 || m[0] != tc.m:
		t.Errorf("ERROR: %q: invalid match: %q (expected %q)",
			tc.s, m, tc.m)
	default:
		t.Logf("%q: success (%q)", tc.s, m)
	}
}
