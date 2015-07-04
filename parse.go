package main

import (
	"fmt"
	"./lex"
)

type Expression struct {
	input string
	Tokens []lex.Token
}

func (exp Expression) String() string {
	return exp.input
}

func (exp *Expression) tokenize() {
	l := lex.Lex(exp.input)

	for token := range l.Tokens {
		exp.Tokens = append(exp.Tokens, token)
	}
}

func newExpression(s string) *Expression {
	exp := &Expression{input: s}

	exp.tokenize()

	return exp
}

// func reciever(c chan lex.Token) {
//     for recievedMsg := range c {
//         fmt.Println("test", recievedMsg)
//     }
// }

func main() {
	x := newExpression("21+4-(xsy*7) ^3")
	fmt.Println(x)
	fmt.Println(x.Tokens)
	// for _,token := range x.Tokens {
	// 	fmt.Println(token)
	// }
}
