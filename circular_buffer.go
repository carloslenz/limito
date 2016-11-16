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

func (b *circularBuffer) add(c chan struct{}) error {
	n := len(b.buf)
	if (b.read == b.write && b.read != 0) || b.write == n {
		return errCircularBufferFull
	}
	next := b.write + 1
	if next > n {
		next %= n
	}
	b.buf[b.write%n] = c
	b.write = next
	return nil
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
