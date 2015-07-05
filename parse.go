package main

import (
	"./lex"
	"fmt"
)

type strategy func([]*Expression) Expression

type Operator struct {
	str    string
	handle strategy
}

var Operators = map[string]Operator{
	"+":  Operator{str: "+"},
	"-":  Operator{str: "-"},
	"**": Operator{str: "**"},
	"^^": Operator{str: "^^"},
	"/":  Operator{str: "/"},
	"-:": Operator{str: "-:"},
	// "-":    tokenOperator,
	// "*":    tokenOperator,
	// "**":   tokenOperator,
	// "^^":   tokenOperator,
	// "/":    tokenOperator,
	// "-:":   tokenOperator,
	// "sum":  tokenOperator,
	// "prod": tokenOperator,

	// "=":  tokenRelation,
	// "!=": tokenRelation,
	// "<":  tokenRelation,
	// ">":  tokenRelation,
	// "<=": tokenRelation,
	// ">=": tokenRelation,

	// "and": tokenLogical,
	// "or":  tokenLogical,
	// "if":  tokenLogical,
	// "iff": tokenLogical,
	// "not": tokenLogical,
	// "=>":  tokenLogical,

	// "sin":  tokenFunction,
	// "cos":  tokenFunction,
	// "tan":  tokenFunction,
	// "csc":  tokenFunction,
	// "sec":  tokenFunction,
	// "cot":  tokenFunction,
	// "sinh": tokenFunction,
	// "cosh": tokenFunction,
	// "tanh": tokenFunction,
	// "log":  tokenFunction,
	// "ln":   tokenFunction,
	// "det":  tokenFunction,
	// "dim":  tokenFunction,
	// "lim":  tokenFunction,
	// "mod":  tokenFunction,
	// "gcd":  tokenFunction,
	// "lcm":  tokenFunction,
	// "min":  tokenFunction,
	// "max":  tokenFunction,
}

// type tokenTransformer func([]lex.Token) lex.Token

// A pattern is a specification of a single kind of
// special token pattern. When a pattern matches, it
// is applied to reduce the size of the linked list
// representation.
type Pattern struct {
	tokens []lex.Token
}

// type Token struct {
// 	value string
// 	typ   tokenType
// 	pos   int
// }

// An expression is an abstract syntax tree representing
// an input string
type Expression struct {
	// String representation of the expression
	raw_input  string
	Tokens     []lex.Token
	Components []lex.Component

	// This describes the tree structure. Each expression
	// forms a node in a tree. The node's children are the
	// inputs to it's operator - there can arbitrarily many
	// operator inputs, so long as the operator of that type
	// understands what to do with the inputs
	operator       Operator
	operatorInputs []*Expression

	// This describes all of the possible adornm0ents to an
	// expression. It forms the state of any given expression
	// This will all serialize into operatorInputs eventually
	superscript *Expression
	subscript   *Expression
	arguments   *Expression
}

func (exp Expression) String() string {
	return exp.raw_input
}

func (exp *Expression) tokenize() {
	l := lex.Lex(exp.raw_input)

	for token := range l.Tokens {
		exp.Tokens = append(exp.Tokens, token)
	}
}

func (exp *Expression) reduceTokens() {
	l := lex.LexTokens(exp.Tokens)

	for component := range l.Components {
		exp.Components = append(exp.Components, component)
	}
}

func newExpression(s string) *Expression {
	exp := &Expression{raw_input: s}

	exp.tokenize()
	exp.reduceTokens()
	// exp.expandTokenTree()

	// signedNumber := Pattern{
	// 	tokens: []lex.Token{
	// 		lex.MakeToken("("),
	// 		lex.MakeToken("("),
	// 	},
	// }

	// fmt.Println(signedNumber)

	return exp
}

// func reciever(c chan lex.Token) {
//     for recievedMsg := range c {
//         fmt.Println("test", recievedMsg)
//     }
// }

func main() {
	x := newExpression("21+4-(xsy*7) ^3")
	// x := newExpression("sum_(i=1)^n i^3=((n(n+1))/2)^2")
	fmt.Println(x)
	fmt.Println(x.Tokens)
	fmt.Println(x.Components)

	// for _,token := range x.Tokens {
	// 	fmt.Println(token)
	// }
}
