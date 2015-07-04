package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	value string
	typ tokenType
	pos int
}

func (i Token) String() string {

	typ := TypeStrings[i.typ]
	return fmt.Sprintf("%[2]s: \"%[1]s\" ", i.value, typ)
}

type Lexer struct {
	input string
	Tokens chan Token
	state stateFn
	start int
	pos int
	width int
}

func (l Lexer) String() string {
	return fmt.Sprintf("start: %[1]d  pos: %[2]d", l.start, l.pos)
}

func (l *Lexer) run() {
	l.state = lexStart

	for l.state != nil {
		l.state = l.state(l)
	}
	close(l.Tokens)
}

// emit a token back to the handler
func (l *Lexer) emit(typ tokenType) Token {
	t := l.getToken(typ)
	l.Tokens <- t
	l.start = l.pos

	return t
}

func (l *Lexer) getToken(typ tokenType) Token {
	return Token{
		value: l.input[l.start:l.pos],
		typ: typ,
		pos: l.start,
	}
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *Lexer) currentRunes() []rune {
	var runes []rune

	maxWidth := l.pos - l.start

	for width := 0; width < maxWidth; {
		r, w := utf8.DecodeRuneInString(l.input[l.start + width:l.pos])

		width += w
		runes = append(runes, r)
	}

	return runes
}

func (l *Lexer) current() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) peek(step int) rune {
	r := l.next()
	l.backup()
	return r
}


func Lex(input string) *Lexer {
	l := &Lexer{
		start: 0,
		pos: 0,
		input: input,
		Tokens: make(chan Token),
	}

	go l.run()

	return l
}

func print(args ...rune) {
	for _,v := range args {
		fmt.Printf("%q ", v)
	}
	fmt.Printf("\n")
}

func prints(args ...string) {
	for _,v := range args {
		fmt.Printf("%s ", v)
	}
	fmt.Printf("\n")
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isGrouping(r rune) bool {
	return strings.ContainsRune(Grouping, r)
}

func isSpecial(r rune) bool {
	return strings.ContainsRune(Special, r)
}

func isAlpha(s string) bool {
	for _,v := range s {
		if !unicode.IsLetter(v) {
			return false
		}
	}
	return true
}

func isDigit(s string) bool {
	for _,v := range s {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}

func isWhitespace(r rune) bool {
	return strings.ContainsRune(Whitespace, r)
}

func isReserved(s string) bool {
	if _, ok := TokenTypes[s]; ok {
	    return true
	}
	return false
}

func getTokenType(s string) tokenType {
	if val, ok := TokenTypes[s]; ok {
	    return val
	}
	return tokenError
}
