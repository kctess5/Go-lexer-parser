package lex

const Alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const Digit = "1234567890"
const Grouping = "(){}[]"
const GroupingEnd = ")]}"
const GroupingStart = "({["
const Special = "+-*^_/:=!<>\\'\""
const Whitespace = " "
const Decorator = "_^"

type tokenType int

const eof = -1
const RuneNull = '\uFFFD'

const (
	tokenNumber tokenType = iota
	tokenOperator
	tokenEOF
	tokenKeyword
	tokenError
	tokenComment
	tokenIdentifier
	tokenRelation
	tokenMisc
	tokenLogical
	tokenGrouping
	tokenFunction
	tokenVariable
	tokenDecorator
	tokenWhitespace
)

// Mapping between computer and human friendly types
var TypeStrings = map[tokenType]string{
	tokenWhitespace: "Whitespace",
	tokenNumber:     "Number",
	tokenOperator:   "Operator",
	tokenEOF:        "EOF",
	tokenKeyword:    "Keyword",
	tokenError:      "Error",
	tokenComment:    "Comment",
	tokenIdentifier: "Identifier",
	tokenRelation:   "Relation",
	tokenMisc:       "Misc",
	tokenLogical:    "Logical",
	tokenGrouping:   "Grouping",
	tokenFunction:   "Function",
	tokenVariable:   "Variable",
	tokenDecorator:  "Decorator",
}

// Mapping between special tokens and types
var TokenTypes = map[string]tokenType{
	" ": tokenWhitespace,

	"^": tokenDecorator,
	"_": tokenDecorator,

	"+":    tokenOperator,
	"-":    tokenOperator,
	"*":    tokenOperator,
	"**":   tokenOperator,
	"^^":   tokenOperator,
	"/":    tokenOperator,
	"-:":   tokenOperator,
	"sum":  tokenOperator,
	"prod": tokenOperator,

	"=":  tokenRelation,
	"!=": tokenRelation,
	"<":  tokenRelation,
	">":  tokenRelation,
	"<=": tokenRelation,
	">=": tokenRelation,

	"and": tokenLogical,
	"or":  tokenLogical,
	"if":  tokenLogical,
	"iff": tokenLogical,
	"not": tokenLogical,
	"=>":  tokenLogical,

	"{": tokenGrouping,
	"}": tokenGrouping,
	"[": tokenGrouping,
	"]": tokenGrouping,
	"(": tokenGrouping,
	")": tokenGrouping,

	"sin":  tokenFunction,
	"cos":  tokenFunction,
	"tan":  tokenFunction,
	"csc":  tokenFunction,
	"sec":  tokenFunction,
	"cot":  tokenFunction,
	"sinh": tokenFunction,
	"cosh": tokenFunction,
	"tanh": tokenFunction,
	"log":  tokenFunction,
	"ln":   tokenFunction,
	"det":  tokenFunction,
	"dim":  tokenFunction,
	"lim":  tokenFunction,
	"mod":  tokenFunction,
	"gcd":  tokenFunction,
	"lcm":  tokenFunction,
	"min":  tokenFunction,
	"max":  tokenFunction,
}
