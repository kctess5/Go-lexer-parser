package lex

import (
	"fmt"
)

type stateFn func(*Lexer) stateFn

func lexAlphaNumeric(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()

			switch w := l.current(); {
			case isDigit(w):
				l.emit(tokenNumber)
				return lexStart
			case isReserved(w):
				l.emit(getTokenType(w))
				return lexStart
			default:
				l.emit(tokenVariable)
				return lexStart
			}

			break Loop
		}
	}
	return lexStart
}

func lexGrouping(l *Lexer) stateFn {
	l.next()
	l.emit(tokenGrouping)
	return lexStart
}

func lexSpecial(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isSpecial(r):
			// absorb.
		default:
			l.backup()

			switch w := l.current(); {
			case isReserved(w):
				l.emit(getTokenType(w))
				return lexStart
			default:
				return l.error("Unexpected Special Characters")
			}

			break Loop
		}
	}
	return lexStart
}

func lexStart(l *Lexer) stateFn {
	switch n := l.next(); {
	case isAlphaNumeric(n):
		l.backup()
		return lexAlphaNumeric
	case isGrouping(n):
		l.backup()
		return lexGrouping
	case isSpecial(n):
		l.backup()
		return lexSpecial
	case isWhitespace(n):
		l.emit(tokenWhitespace)
		return lexStart
	case n == eof:
		l.emit(tokenEOF)
		return nil
	default:
		return l.error("Unimplemented Error")
	}
}

func (l *Lexer) error(err string) stateFn {
	t := l.getToken(tokenError)

	fmt.Println("Error Parsing Token", t)
	fmt.Println("Failed with error:", err)

	return nil
}