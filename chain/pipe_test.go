package chain

import (
	"testing"
	"time"
	"fmt"
	"sync/atomic"
)

func TestSingleProducerOnSingleConsumer(t *testing.T) {
	total := 4

	f := &fakeWidgetHandler{}
	p := NewProducer(total)
	c := NewBaseConsumerWithHandler(0, f)
	w := NewPipe(p, c)

	w.Run()
	w.WaitUntilDone()

	tt := atomic.LoadInt64(&f.total)
	if int(tt) != total {
		t.Errorf("Unexpected handled size, expected %d but has %d", total, tt)
	}
}

func TestSingleProducerOnMultipleBaseConsumers(t *testing.T) {
	total := 4
	consumerPoolSize := 4

	f := &fakeWidgetHandler{}
	p := NewProducer(total)
	b := NewBaseConsumerWithHandler(0, f)
	c := NewConsumerPool(consumerPoolSize, f, b)
	w := NewPipe(p, c)

	w.Run()
	w.WaitUntilDone()

	tt := atomic.LoadInt64(&f.total)
	if int(tt) != total {
		t.Errorf("Unexpected handled size, expected %d but has %d", total, tt)
	}
}

func TestSingleProducerOnMultipleDelayedConsumers(t *testing.T) {
	total := 4
	consumerPoolSize := 4

	f := &fakeWidgetHandler{}
	p := NewProducer(total)
	b := NewDelayedConsumer(0, time.Millisecond * 100, f)
	c := NewConsumerPool(consumerPoolSize, f, b)
	w := NewPipe(p, c)

	w.Run()
	w.WaitUntilDone()

	tt := atomic.LoadInt64(&f.total)
	if int(tt) != total {
		t.Errorf("Unexpected handled size, expected %d but has %d", total, tt)
	}
}

func TestSingleProducerOnDelayedConsumers(t *testing.T) {
	total := 4

	f := &fakeWidgetHandler{}
	d := NewDelayedConsumer(0, time.Second, f)
	p := NewProducer(total)
	w := NewPipe(p, d)

	w.Run()
	w.WaitUntilDone()

	tt := atomic.LoadInt64(&f.total)
	if int(tt) != total {
		t.Errorf("Unexpected handled size, expected %d but has %d", total, tt)
	}
}

type fakeOverRunHandler struct{
	totalHandled *int64
}
func (f *fakeOverRunHandler) Handle(w *widget){
	atomic.AddInt64(f.totalHandled, 1)
}

func TestNonBlockingProducerOnSingleConsumers(t *testing.T) {
	total := 1000

	f := &fakeWidgetHandler{}
	r := &fakeOverRunHandler{&f.total}
	p := NewNonBlockingProducer(total, r)
	c := NewBaseConsumerWithHandler(0, f)
	w := NewPipe(p, c)

	w.Run()
	w.WaitUntilDone()

	tt := atomic.LoadInt64(&f.total)
	if int(tt) != total {
		t.Errorf("Unexpected handled size, expected %d but has %d", total, tt)
	}

}

func TestNonBlockingProducerWithBufferOverRunOnSingleConsumers(t *testing.T) {
	total := 1000

	f := &fakeWidgetHandler{}
	ov := NewBufferedHandler(total)
	p := NewNonBlockingProducer(total, ov)
	c := NewBaseConsumerWithHandler(0, f)
	w := NewPipe(p, c)

	w.Run()
	w.WaitUntilDone()

	tt := atomic.LoadInt64(&f.total)
	if int(tt) != total {
		t.Errorf("Unexpected handled size, expected %d but has %d", total, tt)
	}

}

type fakeWidgetHandler struct {
	total int64
}

func (f *fakeWidgetHandler) Handle(msg string) {
	fmt.Println(msg)
	atomic.AddInt64(&f.total, 1)
}
