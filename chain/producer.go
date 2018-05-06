package chain

type baseProducer struct {
	total int
	index int64
}

func NewProducer(t int) *baseProducer {
	return &baseProducer{
		total: t,
	}
}

func (p *baseProducer) Produce(data chan *widget) {
	defer close(data)
	for i := 0; i < p.total; i++ {
		data <- newWidget(i)
	}
}

// OverrunHandler handle widgets rejected on insertion
type OverrunHandler interface {
	Handle(w *widget)
}

// BufferedOverrunHandler stores rejected widgets on buffer
type BufferedOverrunHandler interface {
	Buffer()[]*widget
}

type nonBlockingProducer struct {
	total      int
	index      int64
	onRejected OverrunHandler
}

// NewNonBlockingProducer produces widgets in non blocking way, on insertion block overrun handler
// is applied
func NewNonBlockingProducer(t int, r OverrunHandler) *nonBlockingProducer {
	return &nonBlockingProducer{
		total:      t,
		onRejected: r,
	}
}

func (p *nonBlockingProducer) Produce(data chan *widget) {
	defer close(data)
	for i := 0; i < p.total; i++ {
		w := newWidget(i)

		select {
		case data <- w:
		default:
			p.onRejected.Handle(w)
		}
	}

	if v, ok := p.onRejected.(BufferedOverrunHandler); ok {
		for _, w := range v.Buffer() {
			data <- w
		}
	}
}

type discardHandler struct {
}

// NewDiscardHandler NOP widget handler
func NewDiscardHandler() *discardHandler{
	return &discardHandler{}
}

func (d *discardHandler)Handle(w *widget) {
}

type bufferedHandler struct{
	size int
	buffer []*widget
}

// NewBufferedHandler buffer rejected widgets on a capped list
func NewBufferedHandler(size int) *bufferedHandler {
	return &bufferedHandler{
		size: size,
		buffer: make([]*widget, 0),
	}
}

func (b *bufferedHandler) Handle(w *widget) {
	b.buffer = append(b.buffer, w)
	if len(b.buffer) > b.size {
		b.buffer = b.buffer[1:]
	}
}

func (b *bufferedHandler) Buffer()[]*widget {
	return b.buffer
}