package glob_test

import (
	"testing"

	"darvaza.org/core"

	"darvaza.org/sidecar/pkg/glob"
)

var _ core.TestCase = captureTestCase{}

// captureMatchCase is a single fixture run through a compiled glob's
// Capture, with the groups it is expected to produce.
type captureMatchCase struct {
	fixture  string
	captures []string
	match    bool
}

func newCaptureMatchCase(fixture string, match bool,
	captures ...string) captureMatchCase {
	// Capture returns a non-nil empty slice for a match with no groups,
	// so keep the expected value non-nil to match under reflect.DeepEqual.
	return captureMatchCase{
		fixture:  fixture,
		captures: core.S(captures...),
		match:    match,
	}
}

// captureTestCase compiles a glob pattern and runs a set of fixtures
// through Capture, checking the captured groups of each.
type captureTestCase struct {
	pattern    string
	separators []rune
	matches    []captureMatchCase
}

func newCaptureTestCase(pattern string, separators []rune,
	matches ...captureMatchCase) captureTestCase {
	return captureTestCase{
		pattern:    pattern,
		separators: separators,
		matches:    matches,
	}
}

// newDotCaptureTestCase builds a captureTestCase using '.' as the only
// separator.
func newDotCaptureTestCase(pattern string,
	matches ...captureMatchCase) captureTestCase {
	return newCaptureTestCase(pattern, []rune{'.'}, matches...)
}

func (tc captureTestCase) Name() string {
	return tc.pattern
}

func (tc captureTestCase) Test(t *testing.T) {
	t.Helper()

	g, err := glob.Compile(tc.pattern, tc.separators...)
	if !core.AssertNoError(t, err, "compile %q", tc.pattern) {
		return
	}

	for _, mc := range tc.matches {
		tc.testMatch(t, g, mc)
	}
}

func (tc captureTestCase) testMatch(t *testing.T, g *glob.Glob,
	mc captureMatchCase) {
	t.Helper()

	captured, ok := g.Capture(mc.fixture)
	if !mc.match {
		core.AssertFalse(t, ok, "%q: match %q", tc.pattern, mc.fixture)
		return
	}

	if core.AssertTrue(t, ok, "%q: match %q", tc.pattern, mc.fixture) {
		core.AssertSliceEqual(t, mc.captures, captured,
			"%q: captures of %q", tc.pattern, mc.fixture)
	}
}

func captureTestCases() []captureTestCase {
	return []captureTestCase{
		newDotCaptureTestCase("*.local",
			newCaptureMatchCase("me.local", true),
			newCaptureMatchCase("local.you", false),
			newCaptureMatchCase("me.local.you", false),
			newCaptureMatchCase("me.you.local", false),
		),
		newDotCaptureTestCase("(*).jpi.io",
			newCaptureMatchCase("www.jpi.io", true, "www"),
		),
	}
}

func TestCapture(t *testing.T) {
	core.RunTestCases(t, captureTestCases())
}
