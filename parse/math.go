package parse

/////////////////////// Grammar Rules ////////////////////////

/*
	These are strictly defined grammar rules that describe all
	allowable configurations of elements within an input string.
	Each rule itself is a parser. The rule can be considered totally
	fulfulled if the parser returns (true, "")

	'mathExpression' is recursively defined through its use of
	'mathComponent.' Because of this recursive call, these rules must
	all be wrapped inside of func's. If one were to simply do:
		var mathExpression = And(...)
	then the Go compiler would yell at you for a typechecking loop!
	In our case, this recursion is not a problem because of the
	or combinator within 'mathComponent' - the result is a tree, where
	all of the mathComponents at the leaves of the tree are simply
	'numbers' and there is no infinite recursive loop.
*/

// any numeric symbol
func digit(s string, n ...string) (bool, string, *Cst) {
	return Or(
		Is("0"), Is("1"), Is("2"), Is("3"), Is("4"),
		Is("5"), Is("6"), Is("7"), Is("8"), Is("9"),
	)(s, "digit")
}

// only sign decorators for numbers
func sign(s string, n ...string) (bool, string, *Cst) {
	return Or(Is("+"), Is("-"))(s, "sign")
}

// simple math operations
func operator(s string, n ...string) (bool, string, *Cst) {
	return Or(
		Is("*"), Is("/"), Is("+"), Is("-"), Is("^"),
	)(s, "operator")
}

/*
	A signed or unsigned number.
	If signed, then it must be surrounded by parentheses.
		e.g. 1337 or (-1337) or (+1337)
	Unsigned numbers can use e notation:
		e.g. 1e4
*/
func number(s string, n ...string) (bool, string, *Cst) {
	return Or(
		And(
			OneOrMore(digit),
			Is("e"),
			OneOrMore(digit),
		),
		OneOrMore(digit),
		And(
			Is("("),
			Optional(sign),
			OneOrMore(digit),
			Is(")"),
		),
	)(s, "number")
}

// A mathComponent is any thing that can go into an expression
// i.e. anything that can be operated upon by an operator
func mathComponent(s string, n ...string) (bool, string, *Cst) {
	return Or(
		number,
		And(Is("("), MathExpression, Is(")")),
	)(s, "mathComponent")
}

// An expression is a string of components separated by operators
func MathExpression(s string, n ...string) (bool, string, *Cst) {
	return And(
		mathComponent,
		Many(And(operator, mathComponent)),
	)(s, "MathExpression")
}
