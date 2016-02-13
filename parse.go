package parse

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"
)

/*
	Utils
*/

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

func RmWhiteSpace(s string) string {
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

type Parser func(*Lexer, ...string) (bool, *Cst)
type ParserCombinator func(...Parser) Parser

type Lexer struct {
	input    string
	position int
}

var UpCounter int
var DownCounter int

func (l *Lexer) advance(n int) {
	if n > 0 {
		UpCounter += n
	}
	l.position += n
}

func (l *Lexer) scanTo(n int) {
	if l.position-n > 0 {
		DownCounter += l.position - n
	}
	l.position = n
}

func (l *Lexer) pos() int {
	return l.position
}

func (l *Lexer) peek(n int) string {
	return l.input[l.position : l.position+n]
}

func (l *Lexer) left() int {
	return len(l.input) - l.position
}

func (l *Lexer) Done() bool {
	return l.left() == 0
}

func (l *Lexer) remainder() string {
	if l.left() > 0 {
		return l.peek(l.left())
	} else {
		return ""
	}
}

func (l *Lexer) peekNextRune() (rune, int) {
	if l.left() >= 5 {
		return utf8.DecodeRuneInString(l.peek(5))
	} else {
		return utf8.DecodeRuneInString(l.remainder())
	}
}

func (l *Lexer) Reset() {
	l.position = 0
}

func NewLexer(s string) *Lexer {
	return &Lexer{s, 0}
}

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
func Is(literal string) Parser {
	return func(l *Lexer, n ...string) (bool, *Cst) {
		// string is too short, fail
		if len(literal) > l.left() {
			return false, nil
		}

		inp := l.peek(len(literal)) // substring of correct length

		if inp == literal {
			name := chooseName(n, nameOf(Is))
			node := NewCst(name)
			node.value = literal

			l.advance(len(literal)) // is a match, shorten string

			return true, node
		} else {
			// not a match, fail
			return false, nil
		}
	}
}

// matches if the input string equals the given literal
func Wildcard(except string) Parser {
	return func(l *Lexer, n ...string) (bool, *Cst) {
		name := chooseName(n, nameOf(Wildcard))
		node := NewCst(name)

		r, w := l.peekNextRune()

		// fmt.Println(r, w, strings.ContainsRune(except, r))

		if w == 0 || strings.ContainsRune(except, r) {
			return false, nil
		} else {
			node.value = string(r)
			l.advance(w)

			return true, node
		}
	}
}

// matches if any one of the given parsers match
func Or(parsers ...Parser) Parser {
	return func(l *Lexer, n ...string) (bool, *Cst) {
		name := chooseName(n, nameOf(Or))
		node := NewCst(name)

		for _, parser := range parsers {
			// test the given parser
			start := l.pos()
			matches, child := parser(l)

			if matches {
				// tree.typ = nameOf(parser)
				node.addChild(child)
				// returns once first parser matches, skips rest
				return true, node
			} else {
				l.scanTo(start)
			}
		}
		// if it has not returned by now, none of the given parsers
		// match, fail
		return false, nil
	}
}

// matches if all of the given parsers match
func And(parsers ...Parser) Parser {
	return func(l *Lexer, n ...string) (bool, *Cst) {
		name := chooseName(n, nameOf(And))
		node := NewCst(name)

		var matches bool
		var child *Cst

		for _, parser := range parsers {
			// test each test, sequentially - i.e. the remainder
			// from the first test is given to the second test, etc

			start := l.pos()

			matches, child = parser(l)

			if !matches {
				l.scanTo(start)

				// if any test fails, the whole thing fails
				return false, nil
			} else {
				node.addChild(child)
			}
		}
		// no test failed, match - return the final remainder
		return true, node
	}
}

/* matches [0...] instances of the given parser, removing all
   of the matched text. i.e.

    Many(Is("a"))("aaaaaaa") ==> (true, "")
*/
func Many(parser Parser) Parser {
	return func(l *Lexer, n ...string) (bool, *Cst) {
		name := chooseName(n, nameOf(Many))
		node := NewCst(name)

		var matches bool
		var child *Cst

		matches = true

		// keeps iterating until the given parser no longer matches
		// feed the remainder forward so that it chomps as it goes
		for matches && l.left() > 0 {
			start := l.pos()
			matches, child = parser(l)

			if matches {
				node.addChild(child)
			} else {
				l.scanTo(start)
			}
		}

		// return whatever remains in the string
		return true, node
	}
}

// matches whether or not the given parser activates
func Optional(parser Parser) Parser {
	return func(l *Lexer, n ...string) (bool, *Cst) {
		name := chooseName(n, nameOf(Optional))
		node := NewCst(name)

		start := l.pos()
		matches, child := parser(l)

		if matches {
			node.addChild(child)
			return true, node
		} else {
			l.scanTo(start)
			return true, node
		}
	}
}

// // matches [1...] instances of the given parser
func OneOrMore(parser Parser) Parser {
	return And(
		parser,       // mandatory match
		Many(parser), // 0 or  subsequent matches
	)
}

func StringParser(p Parser) func(string) (bool, *Cst) {
	return func(s string) (bool, *Cst) {
		l := NewLexer(s)

		matches, tree := p(l)

		if matches && l.left() == 0 {
			return true, tree
		} else {
			fmt.Println("Remainder:", l.remainder()) // l.chomp(l.left()))
			return false, nil
		}
	}
}

func Test() {
	// fmt.Println(StringParser(
	// 	Many(Is("a")),
	// )(`aa`))

	// Many(And(literal, or))

	// matches, _ := StringParser(
	// 	Or(
	// 		Many(And(Or(Is("x"), literal), or)),
	// 		Many(And(literal, or)),
	// 	),
	// )(`'0'|'1'|'2'|'3'|x|'5'|'6'|'7'|'8'|'9'|`)

	matches, _ := StringParser(
		expression,
	)(`'0'`)

	fmt.Println(matches)
	// fmt.Println(StringParser(Is("test"))("test"))
	// fmt.Println(StringParser(Is("a"))("b"))
}
