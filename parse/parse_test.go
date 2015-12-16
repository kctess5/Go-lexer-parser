package parse

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func SetUp(p Parser, s string) (bool, *Cst, *lexer) {
	l := newLexer(s)

	matches, tree := p(l)

	return matches, tree, l
}

func TestIs_Exact(t *testing.T) {
	lit := "test"

	matches, _, l := SetUp(Is(lit), lit)

	assert.True(t, matches, "Same input should match")
	assert.Equal(t, l.pos(), len(lit))
}

func TestIs_NoMatch(t *testing.T) {
	lit := "test"

	matches, _, l := SetUp(Is(lit), "a"+lit)

	assert.False(t, matches, "Different prefix shouldn't match")
	assert.Equal(t, l.pos(), 0)
}

//
func TestIs_PrefixMatch(t *testing.T) {
	lit := "test"

	matches, _, l := SetUp(Is(lit), lit+"a")

	assert.True(t, matches, "Same prefix match")
	assert.Equal(t, l.pos(), len(lit))
}

/*
Wildcard("'") edge cases:
	p("'") => False, 0
	p("a") => True, 1
	p("aaaa") =>
*/

func TestWildcard_Empty(t *testing.T) {
	matches, _, l := SetUp(Wildcard("'"), "")

	assert.False(t, matches, "Wildcard doesn't match empty string")
	assert.Equal(t, l.pos(), 0)
}

func TestWildcard_Matches(t *testing.T) {
	matches, _, l := SetUp(Wildcard("'"), "'")

	assert.False(t, matches, "Matching exception fails")
	assert.Equal(t, l.pos(), 0)
}

func TestWildcard_PrefixMatches(t *testing.T) {
	matches, _, l := SetUp(Wildcard("'"), "aa")

	assert.True(t, matches, "Matching exception fails")
	assert.Equal(t, l.pos(), 1)
}

func TestWildcard_Normal(t *testing.T) {
	matches, _, l := SetUp(Wildcard("'"), "a")

	assert.True(t, matches, "Matching exception fails")
	assert.Equal(t, l.pos(), 1)
}

/*
	Many() Edge cases
*/

func TestMany_None(t *testing.T) {
	matches, _, l := SetUp(Many(Is("a")), "")

	assert.True(t, matches, "Many should match none")
	assert.Equal(t, l.pos(), 0)
}

func TestMany_One(t *testing.T) {
	matches, _, l := SetUp(Many(Is("a")), "a")

	assert.True(t, matches, "Many should match one")
	assert.Equal(t, l.pos(), 1)
}

func TestMany_Many(t *testing.T) {
	matches, _, l := SetUp(Many(Is("a")), "aaa")

	assert.True(t, matches, "Many should match many")
	assert.Equal(t, l.pos(), 3)
}

func TestMany_Prefix(t *testing.T) {
	matches, _, l := SetUp(Many(Is("a")), "bbbaaa")

	assert.True(t, matches, "Many should match wrong prefix without advance")
	assert.Equal(t, l.pos(), 0)
}

func TestMany_Postfix(t *testing.T) {
	matches, _, l := SetUp(Many(Is("a")), "aaabbb")

	assert.True(t, matches, "Many should match wrong postfix")
	assert.Equal(t, l.pos(), 3)
}

func TestMany_Wrong(t *testing.T) {
	matches, _, l := SetUp(Many(Is("a")), "b")

	assert.True(t, matches, "Many should match one wrong w/out advancing")
	assert.Equal(t, l.pos(), 0)
}

/*
	Optional() Edge cases
*/

func TestOptional_None(t *testing.T) {
	matches, _, l := SetUp(Optional(Is("a")), "")

	assert.True(t, matches, "Optional should match null")
	assert.Equal(t, l.pos(), 0)
}

func TestOptional_One(t *testing.T) {
	matches, _, l := SetUp(Optional(Is("a")), "a")

	assert.True(t, matches, "Optional should match one")
	assert.Equal(t, l.pos(), 1)
}

func TestOptional_Wrong(t *testing.T) {
	matches, _, l := SetUp(Optional(Is("a")), "b")

	assert.True(t, matches, "Optional should match wrong without advancing")
	assert.Equal(t, l.pos(), 0)
}
