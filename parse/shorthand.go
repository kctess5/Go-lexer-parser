package parse

import (
	"fmt"
	"strings"
	// "reflect"
	// "runtime"
)

///////////////////// Grammar Rules ////////////////////////

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
// expression -> component [ or component ]
// 			  |  component [ and component ]

func and(s string, n ...string) (bool, string, *Cst) {
	return Is("&")(s, "and")
}
func or(s string, n ...string) (bool, string, *Cst) {
	return Is("|")(s, "or")
}
func literal(s string, n ...string) (bool, string, *Cst) {
	return And(
		Is("'"),
		Many(Wildcard("'")),
		Is("'"),
	)(s, "literal")
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
	return OneOrMore(character)(s, "reference")
}
func many(s string, n ...string) (bool, string, *Cst) {
	return And(
		Is("["),
		expression,
		Is("]"),
	)(s, "many")
}
func optional(s string, n ...string) (bool, string, *Cst) {
	return And(
		Is("("),
		expression,
		Is(")"),
	)(s, "optional")
}
func component(s string, n ...string) (bool, string, *Cst) {
	return Or(
		literal,
		reference,
		many,
		optional,
		expression,
	)(s, "component")
}
func expression(s string, n ...string) (bool, string, *Cst) {
	return Or(
		And(component, Many(And(or, component))),
		And(component, Many(And(and, component))),
	)(s, "expression")
}

/////////////////////// Example ////////////////////////

var math = Grammar{
	"digit":    "'0'|'1'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'",
	"sign":     ` '+'|'-' `,
	"operator": " '*'|'/'|'+'|'-'|'^' ",
	"digits":   "digit & [digit] ",
	"number": ` digits 
			  | digits & 'e' & digits
			  | '(' & (sign) & digits & ')'`,
	"component":  "number | '(' & expression & ')'",
	"expression": "component & [operator & component]",
}

func (g Grammar) Rule(s string) string {
	rule := g[s]
	return rmWhiteSpace(rule)
}

func (g Grammar) GetParser(s string) Parser {
	rule := g.Rule(s)

	matches, remainder, tree := expression(rule)

	if matches && len(remainder) == 0 {
		return GenerateParser(tree)
	} else {
		return nil
	}
}

func GenerateParser(tree *Cst) Parser {

	tree.traverse()

	// func() { fmt.Println("test") }

	return Is("a")
}

func (tree *Cst) traverse() {
	// fxn()
	// fmt.Println(tree.typ)
	// fmt.Println(tree.children)
	if tree.children != nil {
		for child := range tree.childen {
			// 		child.traverse(fxn)
		}
	}

}

/////////////////////////// Main ///////////////////////////////

// func isValid(rule Parser, s string) bool {
// 	matches, remainder, _ := rule(s)
// 	return matches && len(remainder) == 0
// }

func rmWhiteSpace(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func TestShorthand() {
	log := fmt.Println

	log(math.GetParser("digit")("a"))

	// log(expression("[digit]"))
	// log(expression(math.Rule("expression")))
	// log(math)

	// log(expression("1111+111"))

}

// func main() {

// many -> ()
// ./

// log(isValid(expression, "1+2+3"))       // true
// log(isValid(expression, "1+2+3+"))      // false
// log(isValid(expression, "1+(1+(1+1))")) // true
// }
