package main

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"unicode/utf8"
)

/*
	Utils
*/

const prefix = "main"
const verbose = true

var shorthandGrammar = Grammar{
	AND:        "true",
	OR:         "true",
	LITERAL:    "true",
	CHARACTER:  "true",
	REFERENCE:  "true",
	MANY:       "true",
	OPTIONAL:   "true",
	COMPONENT:  "true",
	EXPRESSION: "true",
}

func nameOf(i interface{}) string {
	// this is the string of the function that called the function
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()[len(prefix)+1:]
}

func (t Cst) String() string {
	if verbose {
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
	} else if shorthandGrammar.has(t.typ) {
		// if output
		output := t.typ

		if len(t.value) > 0 {
			output += "<" + t.value + ">"
		}

		childRepr := ""
		for i, child := range t.children {
			if i > 0 {
				childRepr += ", "
			}
			childRepr += child.String()
		}

		if len(childRepr) > 0 {
			left, _, _ := Is("(")(childRepr[:1])
			right, _, _ := Is(")")(childRepr[len(childRepr)-1:])

			if !left || !right {
				output += "(" + childRepr + ")"
			} else {
				output += childRepr
			}
		}

		return output
	} else {
		childrenRep := ""

		for _, child := range t.children {
			childRep := child.String()

			if childRep != "" {
				if childrenRep != "" {
					childrenRep += ", "
				}
				childrenRep += childRep
			}
		}

		return childrenRep
	}
}

func chooseName(s []string, dflt string) string {
	if len(s) > 0 {
		return s[0]
	} else {
		return dflt
	}
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
	test     []string
}

func Ast(name string, optionalChildren ...[]*Cst) *Cst {
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
		node := Ast(name)

		r, w := utf8.DecodeRuneInString(s)

		if w == 0 || strings.ContainsRune(except, r) {
			return false, s, nil
		} else {
			node.value = string(r)
			return true, s[w:], node
		}

		// if w > 0

		// // string is too short, fail
		// if len(s) == 0 || strings.Contains(s, string([]rune(s)[0])) {
		// 	return false, s, nil
		// } else {
		// 	node.value = string([]rune(s)[0])
		// 	return true, s[1:0], node
		// }
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
			node := Ast(name)
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
		node := Ast(name)

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
		node := Ast(name)

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
		node := Ast(name)

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
		node := Ast(name)
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
	// node := Ast(many)

	// matches, remainder, tree :=

	// if matches {
	// 	return true, remainder, node
	// } else {
	// 	return false, remainder, nil
	// }

	return And(
		parser,       // mandatory match
		Many(parser), // 0 or  subsequent matches
	)
}

/////////////////////// Grammar Rules ////////////////////////

type Grammar map[string]string

func (g Grammar) has(i string) bool {
	if _, ok := g[i]; ok {
		return true
	} else {
		return false
	}
}

func (g Grammar) Rule(s string) string {
	rule := g[s]
	return rmWhiteSpace(rule)
}

// literal -> ' & * & '
// reference -> character [character]
// many -> [ & expression & ]
// and -> "&"
// or -> "|"
// optional -> ( & expression & )
// component -> literal
// 			  | expression
// 			  | reference
// 			  | many
// 			  | optional
// 			  | { & expression & }
// expression -> component [ or component ]
// 			  |  component [ and component ]

func and(s string, n ...string) (bool, string, *Cst) {
	return Is("&")(s, AND)
}
func or(s string, n ...string) (bool, string, *Cst) {
	return Is("|")(s, OR)
}
func literal(s string, n ...string) (bool, string, *Cst) {
	return And(
		Is("'"),
		Many(Wildcard("'")),
		Is("'"),
	)(s, LITERAL)
}

func literalToParser(tree *Cst, _ Grammar) Parser {
	many := tree.nthChild(1)

	literal := ""

	for _, child := range many.children {
		literal += child.value
	}

	return Is(literal)
}

func character(s string, n ...string) (bool, string, *Cst) {
	return Or(
		Is("q"), Is("w"), Is("e"), Is("r"), Is("t"), Is("y"), Is("u"),
		Is("i"), Is("o"), Is("p"), Is("a"), Is("s"), Is("d"), Is("f"),
		Is("g"), Is("h"), Is("j"), Is("k"), Is("l"), Is("z"), Is("x"),
		Is("c"), Is("v"), Is("b"), Is("n"), Is("m"),
	)(s)
}
func reference(s string, n ...string) (bool, string, *Cst) {
	return OneOrMore(character)(s, REFERENCE)
}
func referenceToParser(tree *Cst, g Grammar) Parser {
	name := tree.nthChild(0).nthChild(0).value

	optionalChars := tree.nthChild(1)

	if len(optionalChars.children) > 0 {
		for _, child := range optionalChars.children {
			name += child.nthChild(0).value
		}
	}

	return g.GetParser(name)
}

func many(s string, n ...string) (bool, string, *Cst) {
	return And(
		Is("["),
		expression,
		Is("]"),
	)(s, MANY)
}
func manyToParser(tree *Cst, g Grammar) Parser {
	child := tree.nthChild(1)
	return Many(expressionToParser(child, g))
}
func optional(s string, n ...string) (bool, string, *Cst) {
	return And(
		Is("("),
		expression,
		Is(")"),
	)(s, OPTIONAL)
}
func optionalToParser(tree *Cst, g Grammar) Parser {
	child := tree.nthChild(1)
	return Optional(expressionToParser(child, g))
}
func component(s string, n ...string) (bool, string, *Cst) {
	return Or(
		literal,
		reference,
		many,
		optional,
		And(Is("{"), expression, Is("}")),
	)(s, COMPONENT)
}
func componentToParser(tree *Cst, g Grammar) Parser {
	child := tree.nthChild(0)
	switch child.typ {
	case "literal":
		return literalToParser(child, g)
	case "expression":
		return expressionToParser(child, g)
	case "reference":
		return referenceToParser(child, g)
	case "many":
		return manyToParser(child, g)
	case "optional":
		return optionalToParser(child, g)
	case "And":
		return expressionToParser(child.nthChild(1), g)
	}

	debug.PrintStack()

	return nil
}

func expression(s string, n ...string) (bool, string, *Cst) {
	return Or(
		And(component, and, Many(And(component, and)), component),
		And(component, or, Many(And(component, or)), component),
		component,
	)(s, EXPRESSION)
}

func expressionToParser(tree *Cst, g Grammar) Parser {

	if tree.nthChild(0).typ == "component" {
		return componentToParser(tree.nthChild(0), g)
	}

	var children = []Parser{
		componentToParser(tree.nthChild(0).nthChild(0), g),
	}

	operator := tree.nthChild(0).nthChild(1).typ

	optionalComponents := tree.nthChild(0).nthChild(2)

	if len(optionalComponents.children) > 0 {
		for _, child := range optionalComponents.children {
			optionalComponent := child.nthChild(0)
			children = append(children,
				componentToParser(optionalComponent, g))
		}
	}

	children = append(children,
		componentToParser(tree.nthChild(0).nthChild(3), g))

	switch operator {
	case "or":
		return Or(children...)
	case "and":
		return And(children...)
	}

	debug.PrintStack()
	return nil
}

const AND = "and"
const OR = "or"
const LITERAL = "literal"
const CHARACTER = "character"
const REFERENCE = "reference"
const MANY = "many"
const OPTIONAL = "optional"
const COMPONENT = "component"
const EXPRESSION = "expression"

var shorthandMap = map[string]Parser{
	AND:        and,
	OR:         or,
	LITERAL:    literal,
	CHARACTER:  character,
	REFERENCE:  reference,
	MANY:       many,
	OPTIONAL:   optional,
	COMPONENT:  component,
	EXPRESSION: expression,
}

func (g Grammar) GetParser(s string) Parser {
	rule := g.Rule(s)

	fmt.Println(s, "=>", rule)

	matches, remainder, tree := expression(rule)

	// fmt.Println(matches, remainder, tree)

	// fmt.Println()

	if matches && len(remainder) == 0 {
		return expressionToParser(tree, g)
	} else {
		fmt.Println("Invalid Expression!")
		return Is("")
	}
}

/////////////////////// Example ////////////////////////

var math = Grammar{
	"digit":    "'0'|'1'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'",
	"sign":     ` '+'|'-' `,
	"operator": " '*'|'/'|'+'|'-'|'^' ",
	"digits":   "digit & [digit] ",
	"number": ` { digits & 'e' & digits } | digits 
			  | { '(' & (sign) & digits & ')' }`,
	"component":  "number | { '(' & expression & ')' }",
	"expression": "component & [{operator & component}]",
}

func main() {
	log := fmt.Println

	log(math.GetParser("digits")("11")) // true
}

// var shorthandGrammar = Grammar{
// 	"operator":   "true",
// 	"number":     "true",
// 	COMPONENT:  "true",
// 	EXPRESSION: "true",
// }

/*
	These are strictly defined shorthandGrammar rules that describe all
	allowable configurations of elements within an input string.

	Each rule itself is a parser. The rule can be considered totally
	fulfulled if the parser returns (true, "")

	'expression' is recursively defined through its use of 'component'
	Because of this recursive call, these rules must all be wrapped
	inside of func's. If one were to simply do:

		var expression = And(...)

	then the Go compiler would yell at you for a typechecking loop!

	In our case, this recursion is not a problem because of the
	or combinator within 'component' - the result is a tree, where
	all of the components at the leaves of the tree are simply
	'numbers' and there is no infinite recursive loop.
*/

// any numeric symbol
// func digit(s string, n ...string) (bool, string, *Cst) {
// 	return Or(
// 		Is("0"), Is("1"), Is("2"), Is("3"), Is("4"),
// 		Is("5"), Is("6"), Is("7"), Is("8"), Is("9"),
// 	)(s)
// }

// // only sign decorators for numbers
// func sign(s string, n ...string) (bool, string, *Cst) {
// 	return Or(Is("+"), Is("-"))(s)
// }

// // simple math operations
// func operatOr(s string, n ...string) (bool, string, *Cst) {
// 	return Or(Is("*"), Is("/"), Is("+"), Is("-"), Is("^"))(s, "operator")
// }

// /*
// 	A signed or unsigned number.
// 	If signed, then it must be surrounded by parentheses.

// 		e.g. 1337 or (-1337) or (+1337)

// 	Unsigned numbers can use e notation:

// 		e.g. 1e4
// */
// func number(s string, n ...string) (bool, string, *Cst) {
// 	return Or(
// 		And(
// 			OneOrMore(digit),
// 			Is("e"),
// 			OneOrMore(digit),
// 		),
// 		OneOrMore(digit),
// 		And(
// 			Is("("),
// 			Optional(sign),
// 			OneOrMore(digit),
// 			Is(")"),
// 		),
// 	)(s, "number")
// }

// // A component is any thing that can go into an expression
// // i.e. anything that can be operated upon by an operator
// func component(s string, n ...string) (bool, string, *Cst) {
// 	return Or(
// 		number,
// 		And(Is("("), expression, Is(")")),
// 	)(s, COMPONENT)
// }

// // An expression is a string of components separated by operators
// func expression(s string, n ...string) (bool, string, *Cst) {
// 	return And(
// 		component,
// 		Many(And(operator, component)),
// 	)(s, EXPRESSION)
// }

/////////////////////////// Main ///////////////////////////////

// func IsValid(rule Parser, s string) bool {
// matches, remainder, _ := rule(s)
// return matches && len(remainder) == 0
// }

// func Test() {
// log := fmt.Println
// log()
// log(expression("1111+111"))
// log(isValid(expression, "1+2+3"))       // true
// log(isValid(expression, "1+2+3+"))      // false
// log(isValid(expression, "1+(1+(1+1))")) // true
// }
