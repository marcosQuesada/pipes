package chain

import (
	"time"
	"fmt"
	"sync"
)

// ConsumerHandler process Consumer widgets
type ConsumerHandler interface {
	Handle(string)
}

type baseConsumer struct {
	id      string
	handler ConsumerHandler
	wg      sync.WaitGroup
}

// NewBaseConsumerWithHandler default widget Consumer
func NewBaseConsumerWithHandler(i int, h ConsumerHandler) *baseConsumer {
	return &baseConsumer{
		id:      fmt.Sprintf("consumer_%d", i),
		handler: h,
	}
}

func (c *baseConsumer) Consume(data chan *widget) {
	defer c.wg.Done()
	c.wg.Add(1)

	for {
		select {
		case w, open := <-data:
			if !open {
				return
			}
			res := c.format(w)

			c.handler.Handle(res)
		}
	}
}

func (c *baseConsumer) WaitUntilDone() {
	c.wg.Wait()
}

func (d *baseConsumer) format(w *widget) string {
	res := fmt.Sprintf("[%s %s] %s", w.label, w.time.Format("15:04:05.000000"), d.id)

	elapsed := time.Now().Sub(w.time)

	return fmt.Sprintf("%s %s", res, elapsed.String())
}

func (d *baseConsumer)Builder(i int) Consumer {
	return NewBaseConsumerWithHandler(i, d.handler)
}

type delayedConsumer struct {
	delay time.Duration
	wg    sync.WaitGroup
	*baseConsumer
}

// NewDelayedConsumer execute task with d duration delay
func NewDelayedConsumer(id int, d time.Duration, h ConsumerHandler) *delayedConsumer {
	return &delayedConsumer{
		delay:        d,
		baseConsumer: NewBaseConsumerWithHandler(id, h),
	}
}

func (d *delayedConsumer) Consume(data chan *widget) {
	defer d.wg.Done()
	d.wg.Add(1)

	for {
		select {
		case w, open := <-data:
			if !open {
				return
			}

			res := d.format(w)

			// sleep until execution
			time.Sleep(d.delay)

			d.handler.Handle(res)
		}
	}
}

func (d *delayedConsumer) WaitUntilDone() {
	d.wg.Wait()
}

func (d *delayedConsumer)Builder(i int) Consumer {
	return NewDelayedConsumer(i, d.delay, d.handler)
}

type consumerPool struct {
	total   int
	wg      sync.WaitGroup
	handler ConsumerHandler
	builder ConsumerBuilder
}

// NewConsumerPool creates a consumer pool applying fan-out pattern
// creates consumer instances from Consumer Builder
func NewConsumerPool(t int, h ConsumerHandler, b ConsumerBuilder) *consumerPool {
	return &consumerPool{
		total:   t,
		handler: h,
		builder: b,
	}
}

func (c *consumerPool) Consume(data chan *widget) {
	c.wg.Add(c.total)

	// spawn pool consumers
	for i := 0; i < c.total; i++ {
		go func(ix int) {
			defer c.wg.Done()

			// create consumer instance
			dc := c.builder.Builder(ix)
			dc.Consume(data)
		}(i)
	}
}

func (c *consumerPool) WaitUntilDone() {
	c.wg.Wait()
}
