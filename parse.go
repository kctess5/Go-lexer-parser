package main

import "fmt"

//////////////////// String Parser Combinators ///////////////////////

/*
	These string parsers take a string as an input, and return true or
	false depending on if that parser matches the input string. The
	parser also returns the remaining parts of the string - if it
	matched, the matching substring is removed, otherwise, the string
	is returned unchanged.
*/

type Parser func(string) (bool, string)
type ParserCombinator func(...Parser) Parser

// matches if the input string equals the given literal
func is(literal string) Parser {
	return func(s string) (bool, string) {

		// string is too short, fail
		if len(literal) > len(s) {
			return false, s
		}

		inp := s[0:len(literal)] // substring of correct length

		if inp == literal {
			// is a match, shorten string
			return true, s[len(literal):]
		} else {
			// not a match, fail
			return false, s
		}
	}
}

// matches if any one of the given parsers match
func or(parsers ...Parser) Parser {
	return func(s string) (bool, string) {
		for _, parser := range parsers {
			// test the given parser
			matches, remainder := parser(s)
			if matches {
				// returns once first parser matches, skips rest
				return true, remainder
			}
		}
		// if it has not returned by now, none of the given parsers
		// match, fail
		return false, s
	}
}

// matches if all of the given parsers match
func and(parsers ...Parser) Parser {
	return func(s string) (bool, string) {
		var matches bool
		var remainder string

		remainder = s

		for _, parser := range parsers {
			// test each test, sequentially - i.e. the remainder
			// from the first test is given to the second test, etc
			matches, remainder = parser(remainder)

			if !matches {
				// if any test fails, the whole thing fails
				return false, s
			}
		}
		// no test failed, match - return the final remainder
		return true, remainder
	}
}

/* matches [0...] instances of the given parser, removing all
   of the matched text. i.e.

    many(is("a"))("aaaaaaa") ==> (true, "")
*/
func many(parser Parser) Parser {
	return func(s string) (bool, string) {
		var matches bool
		var remainder string

		remainder = s
		matches = true

		// keeps iterating until the given parser no longer matches
		// feed the remainder forward so that it chomps as it goes
		for matches && len(remainder) > 0 {
			matches, remainder = parser(remainder)
		}

		// return whatever remains in the string
		return true, remainder
	}
}

// matches [1...] instances of the given parser
func oneOrMore(parser Parser) Parser {
	return and(
		parser,       // mandatory match
		many(parser), // 0 or  subsequent matches
	)
}

// matches whether or not the given parser activates
func optional(parser Parser) Parser {
	return func(s string) (bool, string) {
		matches, remainder := parser(s)

		if matches {
			return true, remainder
		} else {
			return true, s
		}
	}
}

/////////////////////// Grammar Rules ////////////////////////

/*
	These are strictly defined grammar rules that describe all
	allowable configurations of elements within an input string.

	Each rule itself is a parser. The rule can be considered totally
	fulfulled if the parser returns (true, "")

	'expression' is recursively defined through its use of 'component'
	Because of this recursive call, these rules must all be wrapped
	inside of func's. If one were to simply do:

		var expression = and(...)

	then the Go compiler would yell at you for a typechecking loop!

	In our case, this recursion is not a problem because of the
	or combinator within 'component' - the result is a tree, where
	all of the components at the leaves of the tree are simply
	'numbers' and there is no infinite recursive loop.
*/

// any numeric symbol
func digit(s string) (bool, string) {
	return or(
		is("0"), is("1"), is("2"), is("3"), is("4"),
		is("5"), is("6"), is("7"), is("8"), is("9"),
	)(s)
}

// only sign decorators for numbers
func sign(s string) (bool, string) {
	return or(is("+"), is("-"))(s)
}

// simple math operations
func operator(s string) (bool, string) {
	return or(is("*"), is("/"), is("+"), is("-"), is("^"))(s)
}

/*
	A signed or unsigned number.
	If signed, then it must be surrounded by parentheses.

		e.g. 1337 or (-1337) or (+1337)

	Unsigned numbers can use e notation:

		e.g. 1e4
*/
func number(s string) (bool, string) {
	return or(
		and(
			oneOrMore(digit),
			is("e"),
			oneOrMore(digit),
		),
		oneOrMore(digit),
		and(
			is("("),
			optional(sign),
			oneOrMore(digit),
			is(")"),
		),
	)(s)
}

// A component is any thing that can go into an expression
// i.e. anything that can be operated upon by an operator
func component(s string) (bool, string) {
	return or(
		number,
		and(is("("), expression, is(")")),
	)(s)
}

// An expression is a string of components separated by operators
func expression(s string) (bool, string) {
	return and(
		component,
		many(and(operator, component)),
	)(s)
}

/////////////////////////// Main ///////////////////////////////

func isValid(rule Parser, s string) bool {
	matches, remainder := rule(s)
	return matches && len(remainder) == 0
}

func main() {
	log := fmt.Println

	log(isValid(expression, "1+2+3"))       // true
	log(isValid(expression, "1+2+3+"))      // false
	log(isValid(expression, "1+(1+(1+1))")) // true
}
