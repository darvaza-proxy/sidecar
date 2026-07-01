package httpserver_test

import (
	"testing"

	"darvaza.org/core"

	"darvaza.org/sidecar/pkg/sidecar/httpserver"
)

var _ core.TestCase = acmeHTTP01PatternTestCase{}

// acmeHTTP01PatternTestCase verifies the ACME HTTP-01 challenge path matcher
// against a request path and its expected captured token.
type acmeHTTP01PatternTestCase struct {
	name  string
	path  string // URL.Path under test
	token string // expected captured token ("" when none)
	match bool   // whether the path is a valid challenge path
}

func newAcmeHTTP01PatternTestCase(name, path string, match bool,
	token string) acmeHTTP01PatternTestCase {
	return acmeHTTP01PatternTestCase{
		name:  name,
		path:  path,
		token: token,
		match: match,
	}
}

func (tc acmeHTTP01PatternTestCase) Name() string {
	return tc.name
}

func (tc acmeHTTP01PatternTestCase) Test(t *testing.T) {
	t.Helper()

	m, ok := httpserver.AcmeHTTP01Pattern.Capture(tc.path)

	if !tc.match {
		core.AssertFalse(t, ok, "match %q", tc.path)
		return
	}

	if core.AssertTrue(t, ok, "match %q", tc.path) {
		core.AssertSliceEqual(t, core.S(tc.token), m, "token %q", tc.path)
	}
}

func acmeHTTP01PatternTestCases() []acmeHTTP01PatternTestCase {
	return []acmeHTTP01PatternTestCase{
		newAcmeHTTP01PatternTestCase("empty", "", false, ""),
		newAcmeHTTP01PatternTestCase("root", "/", false, ""),
		newAcmeHTTP01PatternTestCase("unrelated path", "/foo", false, ""),
		newAcmeHTTP01PatternTestCase("well-known only", "/.well-known",
			false, ""),
		newAcmeHTTP01PatternTestCase("well-known slash", "/.well-known/",
			false, ""),
		newAcmeHTTP01PatternTestCase("partial challenge",
			"/.well-known/acme-cha", false, ""),
		newAcmeHTTP01PatternTestCase("partial challenge slash",
			"/.well-known/acme-cha/", false, ""),
		newAcmeHTTP01PatternTestCase("challenge without token",
			"/.well-known/acme-challenge", true, ""),
		newAcmeHTTP01PatternTestCase("challenge trailing slash",
			"/.well-known/acme-challenge/", true, ""),
		newAcmeHTTP01PatternTestCase("challenge with token",
			"/.well-known/acme-challenge/foo", true, "foo"),
		newAcmeHTTP01PatternTestCase("challenge extra segment",
			"/.well-known/acme-challenge/foo/bar", true, ""),
		newAcmeHTTP01PatternTestCase("path traversal prefix",
			"../.well-known/acme-challenge/foo", false, ""),
	}
}

func TestAcmeHTTP01Pattern(t *testing.T) {
	core.RunTestCases(t, acmeHTTP01PatternTestCases())
}
