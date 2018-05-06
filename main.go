package main

import (
	"flag"
	"github.com/marcosQuesada/plato/chain"
	"fmt"
	"time"
)

var totalWidgets = flag.Int("n", 4, "-n=<total> set how many widgets are produced.")
var delay = flag.Int("d", 0, "-d=<delay> set consumer execution delay.")
var consumerPoolSize = flag.Int("c", 1, "-c=<poolSize> set consumer pool size.")
var bufferedOverrun = flag.Int("b", 0, "-b=<bufferSize> buffer entries on insertion block.")
var discardOverrun = flag.Bool("x", false, "-x discard entries on insertion block.")

func main() {
	flag.Parse()

	var p chain.Producer
	var c chain.Consumer

	widgetHandler := &stdoutHandler{}
	p = chain.NewProducer(*totalWidgets)
	c = chain.NewBaseConsumerWithHandler(0, widgetHandler)

	if *delay != 0 {
		c = chain.NewDelayedConsumer(0, time.Millisecond * time.Duration(*delay), widgetHandler)
	}

	if *consumerPoolSize != 1 {
		v, ok := c.(chain.ConsumerBuilder)
		if !ok {
			panic("Chained Consumers need to implement ConsumerBuilder interface")
		}

		c = chain.NewConsumerPool(*consumerPoolSize, widgetHandler, v)
	}

	if *bufferedOverrun != 0 {
		overRun := chain.NewBufferedHandler(*bufferedOverrun)
		p = chain.NewNonBlockingProducer(*totalWidgets, overRun)
	}

	if *discardOverrun {
		p = chain.NewNonBlockingProducer(*totalWidgets, chain.NewDiscardHandler())
	}

	w := chain.NewPipe(p, c)
	w.Run()
	w.WaitUntilDone()
}

// stdoutHandler Prompts entries on StdOut
type stdoutHandler struct {
}

func (s *stdoutHandler) Handle(msg string) {
	fmt.Println(msg)
}
