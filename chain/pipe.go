package chain

import (
	"fmt"
	"time"
)

type widget struct {
	label string
	time  time.Time
}

type Producer interface {
	// Produce widgets and insert them on data channel
	Produce(chan *widget)
}

type Consumer interface {
	// Consume widgets from data channel
	Consume(chan *widget)

	// WaitUntilDone Waits until consumers have processed all widgets
	WaitUntilDone()
}

// ConsumerBuilder handles builder creation, required to handle concurrent consumers
type ConsumerBuilder interface {
	Builder(i int) Consumer
}

// pipe implements producer to builder relation
type pipe struct {
	Producer
	Consumer
	data chan *widget
}

func NewPipe(p Producer, c Consumer) *pipe {
	return &pipe{
		Producer: p,
		Consumer: c,
		data:     make(chan *widget),
	}
}

func (wrk *pipe) Run() {
	go wrk.Consume(wrk.data)

	wrk.Produce(wrk.data)
}

// WaitUntilDone wait until all widgets have been processed
func (wrk *pipe) WaitUntilDone(){
	wrk.Consumer.WaitUntilDone()
}

func newWidget(id int) *widget {
	return &widget{
		label: fmt.Sprintf("widget_id_%d", id),
		time:  time.Now(),
	}
}

