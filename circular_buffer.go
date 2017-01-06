package limito

import "errors"

// circularBuffer uses data types suitable for WaitList.
type circularBuffer struct {
	buf         []chan struct{}
	read, write int
}

var (
	errCircularBufferFull  = errors.New("circular buffer full")
	errCircularBufferEmpty = errors.New("circular buffer empty")
)

func newCircularBuffer(n int) circularBuffer {
	buf := make([]chan struct{}, n)
	for i := 0; i < n; i++ {
		buf[i] = make(chan struct{})
	}

	return circularBuffer{
		buf: buf,
	}
}

func (b *circularBuffer) add() (chan struct{}, error) {
	n := len(b.buf)
	if (b.read == b.write && b.read != 0) || b.write == n {
		return nil, errCircularBufferFull
	}
	next := b.write + 1
	if next > n {
		next %= n
	}
	c := b.buf[b.write%n]
	b.write = next
	return c, nil
}

func (b *circularBuffer) remove() (chan struct{}, error) {
	if b.read == 0 && b.write == 0 {
		return nil, errCircularBufferEmpty
	}

	v := b.buf[b.read]
	b.read++
	if b.read == b.write {
		b.read = 0
		b.write = 0
	}
	return v, nil
}
