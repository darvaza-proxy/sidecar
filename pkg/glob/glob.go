// Package glob implements a glob/capture/replace library
// based on github.com/pachyderm/ohmyglob
package glob

import (
	"time"

	"github.com/dlclark/regexp2"
	"github.com/pachyderm/ohmyglob/compiler"
	"github.com/pachyderm/ohmyglob/syntax"

	"darvaza.org/core"
)

// this file is slightly modified from
// https://github.com/pachyderm/ohmyglob/blob/master/glob.go
// distributed under the MIT licence.

const matchTimeout = 10 * time.Second

// Glob represents a compiled glob pattern
type Glob struct {
	r *regexp2.Regexp
}

// Match tells if a given fixture matches the pattern.
func (g *Glob) Match(fixture string) bool {
	ok, err := g.r.MatchString(fixture)
	if err != nil {
		// timeout
		core.PanicWrap(err, "regexp2.Regexp")
	}
	return ok
}

// Capture returns the list of sub-expressions captured by the pattern
// against the given fixture, and also indicates if there was a match
// at all.
func (g *Glob) Capture(fixture string) ([]string, bool) {
	m, err := g.r.FindStringMatch(fixture)
	switch {
	case err != nil:
		// timeout
		core.PanicWrap(err, "regexp2.Regexp")
	case m == nil:
		// no match
		return nil, false
	}

	groups := m.Groups()
	captures := make([]string, len(groups))

	for _, gp := range groups {
		captures = append(captures, gp.Capture.String())
	}

	return captures, true
}

// Compile creates a [Glob] for the given pattern and separators.
//
// The pattern syntax is:
//
//	pattern:
//	    { term }
//
//	term:
//	    `*`         matches any sequence of non-separator characters
//	    `**`        matches any sequence of characters
//	    `?`         matches any single non-separator character
//	    `[` [ `!` ] { character-range } `]`
//	                character class (must be non-empty)
//	    `{` pattern-list `}`
//	                pattern alternatives
//	    c           matches character c (c != `*`, `**`, `?`, `\`, `[`, `{`, `}`)
//	    `\` c       matches character c
//
//	character-range:
//	    c           matches character c (c != `\\`, `-`, `]`)
//	    `\` c       matches character c
//	    lo `-` hi   matches character c for lo <= c <= hi
//
//	pattern-list:
//	    pattern { `,` pattern }
//	                comma-separated (without spaces) patterns
//
//	captures:
//	    `(` { `|` pattern } `)`
//	    `@(` { `|` pattern } `)`
//	                match and capture one of pipe-separated sub-patterns
//	    `*(` { `|` pattern } `)`
//	                match and capture any number of the pipe-separated sub-patterns
//	    `+(` { `|` pattern } `)`
//	                match and capture one or more of the pipe-separated sub-patterns
//	    `?(` { `|` pattern } `)`
//	                match and capture zero or one of the pipe-separated sub-patterns
//	    `!(` { `|` pattern } `)`
//	                match and capture anything except one of the pipe-separated sub-patterns
func Compile(pattern string, separators ...rune) (*Glob, error) {
	tree, err := syntax.Parse(pattern)
	if err != nil {
		return nil, err
	}

	expr, err := compiler.Compile(tree, separators)
	if err != nil {
		return nil, err
	}

	r, err := regexp2.Compile(expr, regexp2.None)
	if err != nil {
		return nil, err
	}

	r.MatchTimeout = matchTimeout
	return &Glob{r: r}, nil
}

// MustCompile does the same as [Compile] but panics if there is
// an error.
func MustCompile(pattern string, separators ...rune) *Glob {
	g, err := Compile(pattern, separators...)
	if err != nil {
		core.PanicWrap(err, "glob.Compile")
	}
	return g
}
