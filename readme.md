# Go play with parser combinators...

This repo is a simple parser combinator implementation in Go.

In [parse.go](./parse/shorthand.go) you will find the real brains of the repo:

```
Is, Wildcard, Or, And, Many, Optional, OneOrMore
```

These may be functionally composed to parse more interesting things. To aid in this process, I used the combinators to create a shorthand for writing parsers. The shorthand may be found in [shorthand.go](./parse/shorthand.go), and is defined as follows:

```
literal -> ' & [*] & '
reference -> character [character]
many -> [ & expression & ]
and -> "&"
or -> "|"
wildcard -> "* & literal"
optional -> ( & expression & )
component -> literal
		   | expression
		   | reference
		   | many
		   | optional
		   | { & expression & }
		   | wildcard
expression -> component [ or component ]
		   |  component [ and component ]
```

This is slightly easier to understand by example. The following describes a parser that can be used to parse math expressions.

```go
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
```

This may be 'compiled' and used as follows:

```go
input := '1+(1+1)'
// create a lexer for controlling the flow of characters
lexer := parse.NewLexer(input)
// create the parser by parsing the shorthand
parser := math.GetParser("expression")
// parse the input, return a concrete syntax tree & a boolean 'matches'
matches, cst := parser(lex)
```

The concrete syntax tree can be further processed to do something useful, such as evaluating the expression.

Run the examples with:

```
$ go run main.go
```

Hats off to [The Orange Duck](http://theorangeduck.com/page/you-could-have-invented-parser-combinators) for inspiration