package linkedlist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDoublyLinkedList(t *testing.T) {
	l := NewDoublyLinkedList()

	assert.Nil(t, l.head)
	assert.Nil(t, l.tail)
	assert.Zero(t, l.size)
}

func TestSize(t *testing.T) {
	l := NewDoublyLinkedList()

	l.size = 600

	assert.EqualValues(t, l.size, l.Size())
}
func TestHead(t *testing.T) {
	l := NewDoublyLinkedList()

	node := new(Node)

	l.head = node

	assert.Equal(t, node, l.Head())
}

func TestTail(t *testing.T) {
	l := NewDoublyLinkedList()

	node := new(Node)

	l.tail = node

	assert.Equal(t, node, l.Tail())
}

func TestAddToBack(t *testing.T) {
	l := NewDoublyLinkedList()

	vals := []int{1, 2, 3, 4}

	for i, v := range vals {
		lastSize := l.Size()

		node := l.AddToBack(fmt.Sprintf("%d", i), v)

		if i == 0 {
			assert.Equal(t, node, l.Head())
		} else {
			assert.NotEqual(t, node, l.Head())
		}

		assert.Equal(t, node, l.Tail())
		assert.Equal(t, v, node.GetVal())
		assert.Equal(t, fmt.Sprintf("%d", i), node.GetKey())
		assert.Equal(t, lastSize+1, l.Size())
	}

	currNode := l.Head()
	for i := 0; currNode.next != nil; i++ {
		assert.Equal(t, vals[i], currNode.GetVal())

		currNode = currNode.next
	}
}

func TestMoveToBack(t *testing.T) {
	l := NewDoublyLinkedList()

	head := l.AddToBack("1", 1)
	tail := l.AddToBack("2", 2)

	l.MoveToBack(tail)

	assert.Equal(t, tail, l.Tail())
	assert.Equal(t, head, l.Head())

	l.MoveToBack(head)

	assert.Equal(t, tail, l.Head())
	assert.Equal(t, head, l.Tail())

	head = tail
	tail = l.AddToBack("3", 3)

	assert.Equal(t, tail, l.Tail())
	assert.Equal(t, head, l.Head())

	middle := tail

	l.AddToBack("4", 4)

	l.MoveToBack(middle)

	assert.Equal(t, middle, l.Tail())
	assert.Equal(t, head, l.Head())

	l = NewDoublyLinkedList()

	assert.Panics(t, func() {
		l.MoveToBack(tail)
	})
}

func TestRemoveHead(t *testing.T) {
	l := NewDoublyLinkedList()

	l.AddToBack("1", 1)
	h2 := l.AddToBack("2", 2)
	h3 := l.AddToBack("3", 3)

	key := l.RemoveHead()

	assert.Equal(t, key, "1")
	assert.Equal(t, h2, l.Head())
	assert.EqualValues(t, 2, l.Size())

	key = l.RemoveHead()

	assert.Equal(t, key, "2")
	assert.Equal(t, h3, l.Head())
	assert.EqualValues(t, 1, l.Size())

	key = l.RemoveHead()
	assert.Equal(t, key, "3")
	assert.EqualValues(t, 0, l.Size())

	assert.Panics(t, func() {
		l.RemoveHead()
	})
}

func TestNodeMethods(t *testing.T) {
	node := Node{key: "key", value: "value", next: nil, prev: nil}

	assert.Equal(t, "key", node.GetKey())
	assert.Equal(t, "value", node.GetVal())

	node.SetVal("changed")

	assert.Equal(t, "changed", node.GetVal())
}
