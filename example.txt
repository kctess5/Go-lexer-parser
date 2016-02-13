package main

import (
	"./parse"
	"fmt"
	// "github.com/davecheney/profile"
	"time"
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

var json = parse.Grammar{
	"object": `  '{}'
			  |{ '{' & members & '}' } `,

	"members": `{ pair & ',' & members }  
				| pair `,

	"pair": ` string & ':' & value `,

	"array": `  '[]'
			 |{ '[' & elements & ']' } `,

	"elements": `{ value & ',' & elements } 
				 | value `,

	"value": ` 'true'
			 | 'false'
			 | 'null'
			 | string
			 | number
			 | object
			 | array `,

	"string": `  '""'
			  |{ '"' & chars & '"' } `,

	"chars": `{ char & chars}  
			  | char `,

	"char": ` *'\"'
			| '\"'
			| '\\'
			| '\/'
			| '\b'
			| '\f'
			| '\n'
			| '\r'
			| '\t'
			|{'\u' & hexa & hexa & hexa & hexa } `,

	"hexa": ` digit 
			| 'a' | 'b' | 'c' | 'd' | 'e' | 'f'  
			| 'A' | 'B' | 'C' | 'D' | 'E' | 'F'  `,

	"digit": " '0'|'1'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9' ",

	"otn": " '1'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9' ",

	"number": `{ int & frac & exp } 
			  |{ int & frac }
			  |{ int & exp }
			  |  int`,

	"int": `{ '-' & otn & digits } 
		   |{ '-' & digit }
		   |{ otn & digits }
		   |  digit `,

	"frac": " '.' & digits ",

	"exp": " e & digits ",

	"digits": ` { digit & digits } | digit `,

	"e": ` 'e+' | 'e-' | 'E+' | 'E-' | 'e' | 'E' `,
}

func IsValid(g parse.Grammar, ruleName string, input string) bool {
	lex := parse.NewLexer(parse.RmWhiteSpace(input))
	parser := g.GetParser(ruleName)
	matches, _ := parser(lex)
	return matches && lex.Done()
}

func timef(p parse.Parser, s string, iters int) {
	lex := parse.NewLexer(s)
	start := time.Now()

	for i := 0; i < iters; i++ {
		lex.Reset()
		_, _ = p(lex)
	}

	elapsed := time.Since(start)

	fmt.Printf("Function took %[1]s for %[2]d iterations\n", elapsed, iters)
}

func main() {
	// defer profile.Start(profile.CPUProfile).Stop()
	log := fmt.Println

	parse.UpCounter = 0
	parse.DownCounter = 0

	input := `{
		"glossary": {
			"title": "example glossary",
			"GlossDiv": {
				"title": "S",
				"GlossList": {
					"GlossEntry": {
						"ID": "SGML",
						"SortAs": "SGML",
						"GlossTerm": "Standard Generalized Markup Language",
						"Acronym": "SGML",
						"Abbrev": "ISO 8879:1986",
						"GlossDef": {
							"para": "A meta-markup language, used to create markup languages such as DocBook.",
							"GlossSeeAlso": ["GML", "XML"]
						},
						"GlossSee": "markup"
					}
				}
			}
		}
	}`

	// this will benchmark the parser over the given number of iterations
	jsonParser := json.GetParser("object")
	timef(jsonParser, input, 1000)

	// test the math and json parsers
	log(IsValid(json, "object", input))                                         // true
	log(IsValid(math, "expression", "1+(1+(1+(1+(1+(1+(1+(1+(1+(1+1)))))))))")) // true
}
