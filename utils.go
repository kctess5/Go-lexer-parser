package utils

import (
	"./lex"
)

type LinkedListNode struct {
	next LinkedListNode
	prev LinkedListNode
	list LinkedList

	value
}

type LinkedList struct {
	first LinkedListNode
	last  LinkedListNode
}
