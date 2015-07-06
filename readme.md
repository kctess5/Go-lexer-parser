# Go play with parser combinators...

One line says it all:

```go
type Parser func(string) (bool, string)
type ParserCombinator func(...Parser) Parser
```

This at the moment will tell you if a string is a valid math expression,
when it grows up, maybe it will be able to generate the abstract syntax
tree as well!

**parse.go** contains one very well commented implementation of a parser
combinator, and a set of grammar that defines and can test simple 
mathematic operations.

```go
log(isValid(expression, "1+2+3"))       // true
log(isValid(expression, "1+2+3+"))      // false
log(isValid(expression, "1+(1+(1+1))")) // true
```

Run with:

```
$ go run parse.go
```

hats off to [The Orange Duck](http://theorangeduck.com/page/you-could-have-invented-parser-combinators) for inspiration