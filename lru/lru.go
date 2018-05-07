package lru

import (
	"fmt"
)

type node struct {
	key   string
	value int
	head  *node
	tail  *node
}

type lru struct {
	head      *node
	tail      *node
	index     map[string]*node
	size      int
	onEvicted func(k string, val int)
}

func NewLRU(s int, e func(k string, val int)) *lru {
	top := &node{key: "head"}
	tail := &node{key: "tail"}
	top.tail = tail
	tail.head = top

	return &lru{
		head:  top,
		tail:  tail,
		index: make(map[string]*node),
		size:  s,
		onEvicted: e,
	}
}

func (l *lru) Put(k string, v int) {
	defer l.purge()

	if n, ok := l.index[k]; ok {
		n.value = v
		l.clean(n)
		l.setAsHead(n)

		return
	}

	h := l.head

	n := &node{
		head:  l.head,
		tail:  h.tail,
		key:   k,
		value: v,
	}

	l.head.tail.head = n
	l.head.tail = n

	l.index[k] = n
}

func (l *lru) Peek(k string) (int, error) {
	n, ok := l.index[k]
	if !ok {
		return 0, fmt.Errorf("Entry %s not found", k)
	}

	return n.value, nil
}

func (l *lru) Get(k string) (int, error) {
	n, ok := l.index[k]
	if !ok {
		return 0, fmt.Errorf("Entry %s not found", k)
	}

	// remove element from chain
	l.clean(n)

	// Peek to top
	l.setAsHead(n)

	return n.value, nil
}

func (l *lru) purge() {
	if len(l.index) > l.size {
		n := l.tail.head

		if l.onEvicted != nil {
			l.onEvicted(n.key, n.value)
		}

		l.clean(n)

		delete(l.index, n.key)

		l.purge()
	}
}

func (l *lru) setAsHead(n *node) error {
	// Peek to top
	h := l.head
	n.head = h
	n.tail = h.tail

	l.head.tail.head = n
	l.head.tail = n

	return nil
}

// remove element from chain
func (l *lru) clean(n *node) {
	h := n.head
	t := n.tail

	h.tail = t
	t.head = h
}

