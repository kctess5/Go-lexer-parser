package parse

import (
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"
)

/*
	Utils
*/

var prettyPrint = Grammar{
	"and":        "true",
	"or":         "true",
	"literal":    "true",
	"character":  "true",
	"reference":  "true",
	"many":       "true",
	"optional":   "true",
	"component":  "true",
	"expression": "true",
}

func chooseName(s []string, dflt string) string {
	if len(s) > 0 {
		return s[0]
	} else {
		return dflt
	}
}

func nameOf(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	prefixOffset := strings.LastIndex(name, "parse.") + len("parse.")

	// this is the string of the function that called the function
	return name[prefixOffset:]
}

func (t Cst) String() string {
	output := t.typ

	if len(t.value) > 0 {
		output += "<" + t.value + ">"
	}

	childRepr := ""
	for _, child := range t.children {
		childRepr += ", " + child.String()
	}

	if len(childRepr) > 0 {
		output += "(" + childRepr[2:] + ")"
	}

	return output
}

func rmWhiteSpace(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

//////////////////// String Parser Combinators ///////////////////////

/*
	These string parsers take a string as an input, and return true or
	false depending on if that parser matches the input string. The
	parser also returns the remaining parts of the string - if it
	matched, the matching substring is removed, otherwise, the string
	is returned unchanged.

	The string parsers also generate an literal syntax tree for any
	matching sequence. This relies on a simple contract:

	A parser may response will be in one of these two forms:

		- Parser -> (true, _, *Cst)
		- Parser -> (false, _, nil)

	Where *Cst is a pointer to a concrete syntax tree describing the
	matched sequence.
*/

type Parser func(string, ...string) (bool, string, *Cst)
type ParserCombinator func(...Parser) Parser

// type ParserResponse (bool, string, *Cst)

type Cst struct {
	typ      string
	children []*Cst
	value    string
}

func NewCst(name string, optionalChildren ...[]*Cst) *Cst {
	newAst := &Cst{
		typ:      name,
		children: []*Cst{},
	}
	if len(optionalChildren) == 1 {
		newAst.children = optionalChildren[0]
	}
	return newAst
}

func (a *Cst) addChild(child *Cst) {
	a.children = append(a.children, child)
}

func (a *Cst) nthChild(n int) *Cst {
	return a.children[n]
}

// matches if the input string equals the given literal
func Wildcard(except string) Parser {
	return func(s string, n ...string) (bool, string, *Cst) {
		name := chooseName(n, nameOf(Wildcard))
		node := NewCst(name)

		r, w := utf8.DecodeRuneInString(s)

		if w == 0 || strings.ContainsRune(except, r) {
			return false, s, nil
		} else {
			node.value = string(r)
			return true, s[w:], node
		}
	}
}

// matches if the input string equals the given literal
func Is(literal string) Parser {
	return func(s string, n ...string) (bool, string, *Cst) {
		// string is too short, fail
		if len(literal) > len(s) {
			return false, s, nil
		}

		inp := s[0:len(literal)] // substring of correct length

		if inp == literal {
			name := chooseName(n, nameOf(Is))
			node := NewCst(name)
			node.value = literal
			// is a match, shorten string
			return true, s[len(literal):], node
		} else {
			// not a match, fail
			return false, s, nil
		}
	}
}

// matches if any one of the given parsers match
func Or(parsers ...Parser) Parser {
	return func(s string, n ...string) (bool, string, *Cst) {
		name := chooseName(n, nameOf(Or))
		node := NewCst(name)

		for _, parser := range parsers {
			// test the given parser
			matches, remainder, child := parser(s)
			if matches {
				// tree.typ = nameOf(parser)
				node.addChild(child)
				// returns once first parser matches, skips rest
				return true, remainder, node
			}
		}
		// if it has not returned by now, none of the given parsers
		// match, fail
		return false, s, nil
	}
}

// matches if all of the given parsers match
func And(parsers ...Parser) Parser {
	return func(s string, n ...string) (bool, string, *Cst) {
		name := chooseName(n, nameOf(And))
		node := NewCst(name)

		var matches bool
		var remainder string
		var child *Cst

		remainder = s

		for _, parser := range parsers {
			// test each test, sequentially - i.e. the remainder
			// from the first test is given to the second test, etc
			matches, remainder, child = parser(remainder)

			if !matches {
				// if any test fails, the whole thing fails
				return false, s, nil
			}

			// tree.typ = nameOf(parser)

			node.addChild(child)
		}
		// no test failed, match - return the final remainder
		return true, remainder, node
	}
}

/* matches [0...] instances of the given parser, removing all
   of the matched text. i.e.

    Many(Is("a"))("aaaaaaa") ==> (true, "")
*/
func Many(parser Parser) Parser {
	return func(s string, n ...string) (bool, string, *Cst) {
		name := chooseName(n, nameOf(Many))
		node := NewCst(name)

		var matches bool
		var remainder string
		var child *Cst

		remainder = s
		matches = true

		// keeps iterating until the given parser no longer matches
		// feed the remainder forward so that it chomps as it goes
		for matches && len(remainder) > 0 {
			matches, remainder, child = parser(remainder)
			if matches {
				node.addChild(child)
			}
		}

		// return whatever remains in the string
		return true, remainder, node
	}
}

// matches whether or not the given parser activates
func Optional(parser Parser) Parser {
	return func(s string, n ...string) (bool, string, *Cst) {
		name := chooseName(n, nameOf(Optional))
		node := NewCst(name)
		matches, remainder, child := parser(s)

		if matches {
			node.addChild(child)
			return true, remainder, node
		} else {
			return true, s, node
		}
	}
}

// matches [1...] instances of the given parser
func OneOrMore(parser Parser) Parser {
	return And(
		parser,       // mandatory match
		Many(parser), // 0 or  subsequent matches
	)
}
