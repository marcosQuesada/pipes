package lru

import "testing"

func TestLRUDevelopmentFlow(t *testing.T) {

	var called = false
	f := func(k string, val int) {
		called = true
	}
	l := NewLRU(3, f)
	l.Put("a", 1)

	if l.head.tail.key != "a" {
		t.Error("Unexpected first element on LRU")
	}

	if l.tail.head.key != "a" {
		t.Error("Unexpected first element on LRU")
	}

	av, err := l.Peek("a")
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}

	if av != 1 {
		t.Errorf("Unexpected get result, expected 1, got %d", av)
	}

	l.Put("b", 2)
	l.Put("c", 3)
	l.Put("d", 4)

	if !called {
		t.Error("LRU has not been purged!")
	}

	if len(l.index) != 3 {
		t.Errorf("Unexpected index size, expected 3 got %d", len(l.index))
	}

	if l.head.tail.key != "d" {
		t.Error("Unexpected first element on LRU")
	}

	if l.tail.head.key != "b" {
		t.Error("Unexpected last element on LRU")
	}

	l.Put("b", 10)

	if l.head.tail.key != "b" {
		t.Errorf("Unexpected element, got %s", l.head.tail.key)
	}

	if l.tail.head.key != "c" {
		t.Error("Unexpected last element on LRU")
	}

	l.Get("c")

	if l.head.tail.key != "c" {
		t.Errorf("Unexpected element, got %s", l.head.tail.key)
	}

}