package main

import (
	"./parse"
	"fmt"
)

/////////////////////// Example ////////////////////////

var math = parse.Grammar{
	"digit":    "'0'|'1'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'",
	"sign":     ` '+'|'-' `,
	"operator": " '*'|'/'|'+'|'-'|'^' ",
	"digits":   "digit & [digit] ",
	"number": ` { digits & 'e' & digits } | digits 
			  | { '(' & (sign) & digits & ')' }`,
	"component":  "number | { '(' & expression & ')' }",
	"expression": "component & [{operator & component}]",
}

func IsValid(rule parse.Parser, s string) bool {
	matches, remainder, _ := rule(s)
	return matches && len(remainder) == 0
}

func main() {
	log := fmt.Println

	log(math.GetParser("expression")("1+1")) // true, '', Cst
	log(parse.MathExpression("1+1"))

	log(IsValid(math.GetParser("expression"), "1+(1+(1+(1+(1+1))))")) // true
}
