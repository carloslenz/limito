package limito

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func ExampleMiddleware() {
	f := Middleware(1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("hello %s", GetMiddlewareID(r.Context()))
	}))

	ctx := SetMiddlewareID(context.Background(), "john_doe")
	f.ServeHTTP(nil, (&http.Request{}).WithContext(ctx))

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(" - %s\n", err)
		}
	}()
	f.ServeHTTP(nil, (&http.Request{}).WithContext(context.Background()))
	// Output: hello john_doe - limito middleware ID: not found
}

func TestMiddleware(t *testing.T) {
	c := make(chan int)
	f := Middleware(1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		c <- 0
		rw.WriteHeader(http.StatusOK)
	}))

	var done sync.WaitGroup
	var resp [4]*httptest.ResponseRecorder
	serve := func(i int, ctx context.Context) {
		rw := &httptest.ResponseRecorder{}
		resp[i] = rw
		done.Add(1)
		defer done.Done()
		// 1st request waits for chan, so the 2nd fails because limit is 1 and they have the same ID.
		f.ServeHTTP(rw, (&http.Request{}).WithContext(ctx))
	}

	ctx := SetMiddlewareID(context.Background(), "test")
	for i := 0; i < 2; i++ {
		go serve(i, ctx)
	}

	// This request is ok because the ID is different.
	ctx2 := SetMiddlewareID(context.Background(), "another")
	go serve(2, ctx2)

	time.Sleep(10 * time.Millisecond) // So goroutines wait on chan.

	<-c // 1st (or 2nd) request can continue.
	<-c // 3rd request can continue.

	done.Wait() // make sure 1st (or 2nd) and 3rd requestes have finished.

	// Checking that the 1st ID can be used again
	go serve(3, ctx)
	<-c // 4th request can continue.

	done.Wait() // make sure 4th request has finished.

	var codes [len(resp)]int
	for i := 0; i < len(resp); i++ {
		codes[i] = resp[i].Code
	}

	if codes[0] > codes[1] {
		// swapping for comparison: due to concurrency there is no guarantee who arrives 1st (and succeeds).
		codes[0], codes[1] = codes[1], codes[0]
	}

	expected := [len(resp)]int{200, 429, 200, 200}
	if codes != expected {
		t.Errorf("Unexpected HTTP status codes, got %#v, expected %#v", codes, expected)
	}
}
