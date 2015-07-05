package lex

import (
	"fmt"
	"strings"
)

type tlStateFxn func(*TokenLexer) tlStateFxn

// func lexAlphaNumeric(tl *TokenLexer) tlStateFxn {
// Loop:
// 	for {
// 		switch r := l.next(); {
// 		case isAlphaNumeric(r):
// 			// absorb.
// 		default:
// 			l.backup()

// 			switch w := l.current(); {
// 			case isDigit(w):
// 				l.emit(tokenNumber)
// 				return lexStart
// 			case isReserved(w):
// 				l.emit(getTokenType(w))
// 				return lexStart
// 			default:
// 				l.emit(tokenVariable)
// 				return lexStart
// 			}

// 			break Loop
// 		}
// 	}
// 	return lexStart
// }

// func lexGrouping(tl *TokenLexer) tlStateFxn {
// 	l.next()
// 	l.emit(tokenGrouping)
// 	return lexStart
// }

// func lexSpecial(tl *TokenLexer) tlStateFxn {
// Loop:
// 	for {
// 		switch r := l.next(); {
// 		case isSpecial(r):
// 			// absorb.
// 		default:
// 			l.backup()

// 			switch w := l.current(); {
// 			case isReserved(w):
// 				l.emit(getTokenType(w))
// 				return lexStart
// 			default:
// 				return l.error("Unexpected Special Characters")
// 			}

// 			break Loop
// 		}
// 	}
// 	return lexStart
// }

func tokenLexPostAtom(tl *TokenLexer) tlStateFxn {
	// switch n := tl.next(); {
	// case isNumber(n):
	// 	tl.backup()
	// 	return tokenLexNumber
	// tl.next()
	// tl.emit(tokenComment)
	return tl.error("Unimplemented Error")
}

func tokenLexNumber(tl *TokenLexer) tlStateFxn {
	tl.next()
	tl.emit()
	return tokenLexPostAtom
}

func tokenLexGrouping(tl *TokenLexer) tlStateFxn {
	tl.next()
	// if belongsTo(GroupingEnd, t.value) {
	// 	tl.
	// }
	// case t.value
	tl.emit()
	return tokenLexPostAtom
}

func tokenLexStart(tl *TokenLexer) tlStateFxn {
	// tl.emit(tokenEOF)
	switch n := tl.next(); {
	case isNumberToken(n):
		tl.backup()
		return tokenLexNumber
	case isGroupingToken(n):
		tl.backup()
		return tokenLexGrouping
	// case isGrouping(n):
	// 	l.backup()
	// 	return lexGrouping
	// case isSpecial(n):
	// 	l.backup()
	// 	return lexSpecial
	// case isWhitespace(n):
	// 	l.emit(tokenWhitespace)
	// 	return lexStart
	// case n == eof:
	// 	l.emit(tokenEOF)
	// 	return nil
	default:
		return tl.error("Unimplemented Error")
	}
	return tl.error("Unimplemented Error")
}

func (tl *TokenLexer) error(err string) tlStateFxn {

	fmt.Println("Error Parsing Tokens")
	fmt.Println("Failed with error:", err)

	return nil
}

func isNumberToken(t *Token) bool {
	return t.typ == tokenNumber
}
func isGroupingToken(t *Token) bool {
	return t.typ == tokenGrouping
}

func belongsTo(s string, test string) bool {
	return strings.Contains(s, test)
}
