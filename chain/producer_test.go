package chain

import "testing"

func TestOverRunBufferedHandlerCapsSizeWhenFull(t *testing.T) {

	b := NewBufferedHandler(1)

	b.Handle(&widget{})
	b.Handle(&widget{})
	if 1 != len(b.Buffer()) {
		t.Errorf("Unexpected Buffer Size, expected 1 but got %d", len(b.Buffer()))
	}
}
