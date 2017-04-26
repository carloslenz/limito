package limito

import (
	"errors"
	"sync/atomic"
)

// Limiter limits the number of concurrent operations.
type Limiter struct {
	closed int32
	c      chan struct{}
}

// NewLimiter returns a concurrent operations limiter. Supports up to n concurrent operations.
func NewLimiter(n int) Limiter {
	return Limiter{
		c: make(chan struct{}, n),
	}
}

var (
	errLimiterFull  = errors.New("limiter full")
	errLimiterEmpty = errors.New("limiter empty")
)

// Enter locks one slot. Error is returned if all slots are in use.
func (l *Limiter) Enter() error {
	select {
	case l.c <- struct{}{}:
		return nil
	default:
		return errLimiterFull
	}
}

// Leave releases one slot. Error is returned if no slot is in use.
func (l *Limiter) Leave() error {
	select {
	case <-l.c:
		return nil
	default:
		return errLimiterEmpty
	}
}

// Close releases resources. Warning: any call to Enter or Leave will panic after Close.
// Multiple Close calls are ok.
func (l *Limiter) Close() {
	if atomic.CompareAndSwapInt32(&l.closed, 0, 1) {
		close(l.c)
	}
}
