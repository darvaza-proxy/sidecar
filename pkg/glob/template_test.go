package glob_test

// cspell:words orld

import (
	"testing"

	"darvaza.org/core"

	"darvaza.org/sidecar/pkg/glob"
)

var (
	_ core.TestCase = templateCompileTestCase{}
	_ core.TestCase = templateReplaceTestCase{}
	_ core.TestCase = templateEqualTestCase{}
)

func mustCompileTemplate(t *testing.T, s string) *glob.Template {
	t.Helper()

	p, err := glob.CompileTemplate(s)
	core.AssertMustNoError(t, err, "compile %q", s)
	return p
}

// templateCompileTestCase compiles a template string and checks the parsed
// structure through its String representation.
type templateCompileTestCase struct {
	input string
	want  string
	ok    bool
}

func newTemplateCompileTestCase(input, want string) templateCompileTestCase {
	return templateCompileTestCase{
		input: input,
		want:  want,
		ok:    true,
	}
}

func newTemplateCompileErrorTestCase(input string) templateCompileTestCase {
	return templateCompileTestCase{
		input: input,
		ok:    false,
	}
}

func (tc templateCompileTestCase) Name() string {
	return tc.input
}

func (tc templateCompileTestCase) Test(t *testing.T) {
	t.Helper()

	p, err := glob.CompileTemplate(tc.input)
	if !tc.ok {
		core.AssertError(t, err, "compile %q", tc.input)
		return
	}

	if core.AssertNoError(t, err, "compile %q", tc.input) {
		core.AssertEqual(t, tc.want, p.String(), "compile %q", tc.input)
	}
}

func templateCompileTestCases() []templateCompileTestCase {
	return []templateCompileTestCase{
		newTemplateCompileTestCase("hello", `{"hello"}`),
		newTemplateCompileTestCase("hello$", `{"hello$"}`),
		newTemplateCompileTestCase("hello$1", `{"hello", 1}`),
		newTemplateCompileErrorTestCase("hello${1"),
		newTemplateCompileTestCase("hello${1}", `{"hello", 1}`),
		newTemplateCompileTestCase("hello${1}world", `{"hello", 1, "world"}`),
		newTemplateCompileErrorTestCase("hello${1w}orld"),
		newTemplateCompileErrorTestCase("hello${1world"),
		newTemplateCompileErrorTestCase("hello${world"),
		newTemplateCompileErrorTestCase("hello${0}world"),
		newTemplateCompileErrorTestCase("hello${-3}world"),
		newTemplateCompileErrorTestCase("hello${-3wo}rld"),
		newTemplateCompileErrorTestCase("hello${-3world"),
		newTemplateCompileTestCase("hello${1}world${2}from${3}",
			`{"hello", 1, "world", 2, "from", 3}`),
		newTemplateCompileTestCase("a$b$c$5gh", `{"a$b$c", 5, "gh"}`),
	}
}

func TestTemplate(t *testing.T) {
	core.RunTestCases(t, templateCompileTestCases())
}

// templateReplaceTestCase compiles a template, applies Replace with the
// given data and checks the produced string.
type templateReplaceTestCase struct {
	input  string
	result string
	data   []string
	ok     bool
}

func newTemplateReplaceTestCase(input, result string,
	data ...string) templateReplaceTestCase {
	return templateReplaceTestCase{
		input:  input,
		result: result,
		data:   data,
		ok:     true,
	}
}

func newTemplateReplaceErrorTestCase(input string,
	data ...string) templateReplaceTestCase {
	return templateReplaceTestCase{
		input: input,
		data:  data,
		ok:    false,
	}
}

func (tc templateReplaceTestCase) Name() string {
	return tc.input
}

func (tc templateReplaceTestCase) Test(t *testing.T) {
	t.Helper()

	p, err := glob.CompileTemplate(tc.input)
	if !core.AssertNoError(t, err, "compile %q", tc.input) {
		return
	}

	r, err := p.Replace(tc.data)
	if !tc.ok {
		core.AssertError(t, err, "replace %q with %q", tc.input, tc.data)
		return
	}

	if core.AssertNoError(t, err, "replace %q with %q", tc.input, tc.data) {
		core.AssertEqual(t, tc.result, r,
			"replace %q with %q", tc.input, tc.data)
	}
}

func templateReplaceTestCases() []templateReplaceTestCase {
	return []templateReplaceTestCase{
		newTemplateReplaceTestCase("foobar", "foobar"),
		newTemplateReplaceTestCase("$1Foe", "oneFoe", "one"),
		newTemplateReplaceTestCase("${1}Foe", "oneFoe", "one"),
		newTemplateReplaceTestCase("${1}Foe", "Foe", "", "two"),
		newTemplateReplaceErrorTestCase("${2}Foe", "one"),
		newTemplateReplaceTestCase("Hello, ${2}", "Hello, world", "one", "world"),
	}
}

func TestReplace(t *testing.T) {
	core.RunTestCases(t, templateReplaceTestCases())
}

// templateEqualTestCase compiles two templates and checks Template.Equal.
type templateEqualTestCase struct {
	name     string
	left     string
	right    string
	rightNil bool
	want     bool
}

func newTemplateEqualTestCase(name, left, right string,
	want bool) templateEqualTestCase {
	return templateEqualTestCase{
		name:  name,
		left:  left,
		right: right,
		want:  want,
	}
}

func newTemplateEqualNilTestCase(name, left string) templateEqualTestCase {
	return templateEqualTestCase{
		name:     name,
		left:     left,
		rightNil: true,
		want:     false,
	}
}

func (tc templateEqualTestCase) Name() string {
	return tc.name
}

func (tc templateEqualTestCase) Test(t *testing.T) {
	t.Helper()

	left := mustCompileTemplate(t, tc.left)

	var right *glob.Template
	if !tc.rightNil {
		right = mustCompileTemplate(t, tc.right)
	}

	core.AssertEqual(t, tc.want, left.Equal(right), "%s", tc.name)
}

func templateEqualTestCases() []templateEqualTestCase {
	return []templateEqualTestCase{
		newTemplateEqualTestCase("identical",
			"hello${1}world", "hello${1}world", true),
		newTemplateEqualTestCase("different length",
			"hello", "hello${1}", false),
		newTemplateEqualTestCase("different literal",
			"hello", "world", false),
		newTemplateEqualTestCase("different index",
			"$1", "$2", false),
		newTemplateEqualNilTestCase("nil right", "hello"),
	}
}

func TestEqual(t *testing.T) {
	core.RunTestCases(t, templateEqualTestCases())
}
