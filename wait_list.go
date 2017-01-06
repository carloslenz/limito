package limito

import (
	"context"
	"fmt"
	"sync"
)

// WaitList holds calls in line until resource is free.
type WaitList struct {
	mutex   sync.Mutex
	waiting circularBuffer
	blocked bool
}

// NewWaitList creates a WaitList for up to max waiting goroutines.
// Beyond that, requests are rejected.
func NewWaitList(max int) WaitList {
	return WaitList{
		waiting: circularBuffer{
			buf: make([]chan struct{}, max),
		},
	}
}

// Wait waits until the resource is available. Returns error if w is full.
func (w *WaitList) Wait(ctx context.Context) error {
	w.mutex.Lock()
	if w.blocked {
		w.mutex.Unlock()
		return nil
	}

	c := make(chan struct{})
	err := w.waiting.add(c)
	w.blocked = true
	w.mutex.Unlock()
	if err != nil {
		return fmt.Errorf("wait list: %s", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c:
		w.mutex.Lock()
		w.waiting.remove()
		w.mutex.Unlock()
		return nil
	}
}

// Release notifies the oldest waiting gorouting to continue, if any.
// No-op if the w is empty.
func (w *WaitList) Release() {
	w.mutex.Lock()

	c, err := w.waiting.remove()
	if err == errCircularBufferEmpty {
		w.blocked = false
		w.mutex.Unlock()
		return
	}
	if err != nil {
		w.mutex.Unlock()
		// Unexpected circularBuffer behavior: bug!
		panic(err)
	}
	close(c)
	w.mutex.Unlock()
}
