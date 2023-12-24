package glob

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Template is the compiled representation of
// a replace pattern
type Template struct {
	parts []templatePart
}

func (p Template) String() string {
	s := make([]string, len(p.parts))
	for i, p := range p.parts {
		s[i] = p.String()
	}

	return fmt.Sprintf("{%s}", strings.Join(s, ", "))
}

// Equal tells if two instances have the same content
func (p Template) Equal(q *Template) bool {
	switch {
	case q == nil:
		return false
	case len(p.parts) != len(q.parts):
		return false
	default:
		for i := range p.parts {
			if !p.parts[i].Equal(q.parts[i]) {
				return false
			}
		}
	}
	return true
}

// Replace fills gaps in the template with elements of the given
// slice
func (p Template) Replace(data []string) (string, error) {
	out := make([]string, len(p.parts))

	for i, part := range p.parts {
		var s string

		if idx := part.index; idx > 0 {
			if idx >= len(data) {
				err := fmt.Errorf("invalid reference to capture %v", idx)
				return "", err
			}
			s = data[idx-1]
		} else {
			s = part.literal
		}

		out[i] = s
	}

	return strings.Join(out, ""), nil
}

func (p *Template) reduce() {
	parts := make([]templatePart, 0, len(p.parts))

	j := 0
	for _, p := range p.parts {
		if p.literal != "" && j > 0 && parts[j-1].literal != "" {
			// merge literal
			parts[j-1].literal += p.literal
		} else {
			// append
			parts = append(parts, p)
			j++
		}
	}

	if len(parts) != len(p.parts) {
		p.parts = parts
	}
}

type templatePart struct {
	literal string
	index   int
}

func (pp templatePart) Equal(qq templatePart) bool {
	return pp.literal == qq.literal && pp.index == qq.index
}

func (pp templatePart) String() string {
	if pp.index > 0 {
		return fmt.Sprintf("%v", pp.index)
	}

	return fmt.Sprintf("%q", pp.literal)
}

// CompileTemplate creates a [Template] for the given template
//
// Syntax is very simple. `${n}` and `$n` to indicate position
// in the captures slice, and `$$` to escape a literal `$`.
//
// Escaping `$` isn't required when followed by anything other
// than a number or a `{`.
//
// Captures are counted starting with 1.
func CompileTemplate(template string) (*Template, error) {
	out := &Template{}
	rest := template

	for len(rest) > 0 {
		var err error

		rest, err = compileTemplatePass(out, template, rest)
		if err != nil {
			return nil, err
		}
	}

	out.reduce()
	return out, nil
}

func compileTemplatePass(out *Template, template, rest string) (string, error) {
	var literal string

	i := strings.Index(rest, "$")
	switch {
	case i < 0:
		// no more positions
		literal = rest
		rest = ""
	case i > 0:
		// position after a literal string
		literal = rest[:i]
		rest = rest[i:]
	case len(rest) == 1:
		// one lone "$"
		literal = rest
		rest = ""
	case rest[1] == '$':
		// escaped dollar, "$$"
		literal = "$"
		rest = rest[2:]
	default:
		// other cases starting with a "$"
		s0, s1, ok := compileTemplateDollar(out, rest)
		if !ok {
			l := len(template) - len(rest)
			// TODO: test if this counts runes or bytes
			p := len(template[:l])

			err := fmt.Errorf("invalid syntax at position %v", p+1)
			return "", err
		}

		literal = s0
		rest = s1
	}

	if literal != "" {
		out.parts = append(out.parts, templatePart{
			literal: literal,
		})
	}

	return rest, nil
}

// revive:disable:cognitive-complexity
// revive:disable:cyclomatic
func compileTemplateDollar(out *Template, str string) (literal, rest string, ok bool) {
	// revive:enable:cognitive-complexity
	// revive:enable:cyclomatic
	var brace bool
	var sNum string
	var num int

	if str[1] == '{' {
		brace = true
		str = str[2:]
	} else {
		str = str[1:]
	}

	i := 0
	for i < len(str) {
		r, l := utf8.DecodeRuneInString(str[i:])
		if !unicode.IsDigit(r) {
			break
		}
		i += l
	}

	sNum = str[:i]
	rest = str[i:]

	switch {
	case len(sNum) > 0:
		// "$<num>" or "${<num>"
		s64, _ := strconv.ParseInt(sNum, 10, 16)
		num = int(s64)
		if num < 1 {
			// invalid index
			return "", "", false
		}
	case brace:
		// invalid "${"
		return "", "", false
	default:
		// `$x``
		return "$", rest, true
	}

	if brace {
		// "${<num>"
		if len(rest) == 0 || rest[0] != '}' {
			// '}' missing
			return "", "", false
		}
		// "${<num>}"
		rest = rest[1:]
	}

	// append index part
	out.parts = append(out.parts, templatePart{
		index: num,
	})

	return "", rest, true
}
