package parse

import (
	"fmt"
	"runtime/debug"
)

const verbose = false

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
	return RmWhiteSpace(rule)
}

var parserCache = map[string]Parser{}

func (g Grammar) GetParser(s string) Parser {

	// if seen before, and finished:
	// 		return memoized value
	// if seen before, and not finished:
	// 		return promise
	// if not seen before,
	// 		return expressionToParser(rule(s))

	if v, ok := parserCache[s]; ok {
		if v == nil {
			return func(l *Lexer, n ...string) (bool, *Cst) {
				// we must defer the access of the map until parser
				// runtime, otherwise recursively defined grammars
				// would not ever finish compiling
				return parserCache[s](l, s)
			}
		} else {
			return func(l *Lexer, n ...string) (bool, *Cst) {
				return v(l, s) // pass the name of the parser
			}
		}
	}

	lex := NewLexer(g.Rule(s))
	matches, tree := expression(lex)

	if verbose {
		fmt.Println("Generating Parser: ", s, "=>", g.Rule(s))
	}

	parserCache[s] = nil

	if matches && lex.left() == 0 {
		parserCache[s] = expressionToParser(tree, g)

		return func(l *Lexer, n ...string) (bool, *Cst) {
			return parserCache[s](l, s) // pass the name of the parser
		}
	} else {
		fmt.Println("Invalid Parser Expression!")
		return nil
	}
}

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

// literal -> ' & * & '
// reference -> character [character]
// many -> [ & expression & ]
// and -> "&"
// or -> "|"
// wildcard -> "* & literal"
// optional -> ( & expression & )
// component -> literal
// 			  | expression
// 			  | reference
// 			  | many
// 			  | optional
// 			  | { & expression & }
// 			  | wildcard
// expression -> component [ or component ]
// 			  |  component [ and component ]

func wildcard(l *Lexer, n ...string) (bool, *Cst) {
	return And(Is(`*`), literal)(l, "wildcard")
}

func wildcardToParser(tree *Cst, _ Grammar) Parser {
	exclusions := tree.nthChild(1).nthChild(1)

	except := ""

	for _, child := range exclusions.children {
		except += child.value
	}

	return Wildcard(except)
}

func and(l *Lexer, n ...string) (bool, *Cst) {
	return Is("&")(l, "and")
}
func or(l *Lexer, n ...string) (bool, *Cst) {
	return Is("|")(l, "or")
}
func character(l *Lexer, n ...string) (bool, *Cst) {
	return Or(
		Is("q"), Is("w"), Is("e"), Is("r"), Is("t"), Is("y"), Is("u"),
		Is("i"), Is("o"), Is("p"), Is("a"), Is("s"), Is("d"), Is("f"),
		Is("g"), Is("h"), Is("j"), Is("k"), Is("l"), Is("z"), Is("x"),
		Is("c"), Is("v"), Is("b"), Is("n"), Is("m"),
	)(l, "character")
}

func literal(l *Lexer, n ...string) (bool, *Cst) {
	return And(
		Is("'"),
		Many(Wildcard("'")),
		Is("'"),
	)(l, "literal")
}

// func literalToParser(tree *Cst, _ Grammar) Parser {
// 	many := tree.nthChild(1)
// 	lit := ""

// 	for _, child := range many.children {
// 		lit += child.value
// 	}
// 	return Is(lit)
// }

func literalToParser(tree *Cst, _ Grammar) Parser {
	many := tree.nthChild(1)

	characterParsers := []Parser{}

	// literal := ""

	for _, child := range many.children {
		characterParsers = append(characterParsers, Is(child.value))
	}
	return And(characterParsers...)
}

func reference(l *Lexer, n ...string) (bool, *Cst) {
	return OneOrMore(character)(l, "reference")
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

func many(l *Lexer, n ...string) (bool, *Cst) {
	return And(
		Is("["),
		expression,
		Is("]"),
	)(l, "many")
}

func manyToParser(tree *Cst, g Grammar) Parser {
	child := tree.nthChild(1)
	return Many(expressionToParser(child, g))
}

func optional(l *Lexer, n ...string) (bool, *Cst) {
	return And(
		Is("("),
		expression,
		Is(")"),
	)(l, "optional")
}

func optionalToParser(tree *Cst, g Grammar) Parser {
	child := tree.nthChild(1)
	return Optional(expressionToParser(child, g))
}

func component(l *Lexer, n ...string) (bool, *Cst) {
	return Or(
		literal,
		reference,
		many,
		optional,
		wildcard,
		And(Is("{"), expression, Is("}")),
	)(l, "component")
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
	case "wildcard":
		return wildcardToParser(child, g)
	case "And":
		return expressionToParser(child.nthChild(1), g)
	}

	fmt.Println("Unexpected Component", child.typ)
	debug.PrintStack()
	return nil
}

func expression(l *Lexer, n ...string) (bool, *Cst) {
	return Or(
		And(component, and, Many(And(component, and)), component),
		And(component, or, Many(And(component, or)), component),
		component,
	)(l, "expression")
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
