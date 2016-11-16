package limito

import (
	"context"
	"fmt"
	"net/http"
	"testing"
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
	var res []int
	f := Middleware(1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		res = append(res, len(res))
		// 1st request waits here, so the 2nd doesn't start thanks to the middleware.
		c <- 0
	}))
	ctx := SetMiddlewareID(context.Background(), "test")
	for i := 0; i < 2; i++ {
		go func() {
			f.ServeHTTP(nil, (&http.Request{}).WithContext(ctx))
		}()
	}

	var expected []int
	check := func(i int, l ...int) {
		expected = append(expected, l...)
		a := fmt.Sprint(res)
		b := fmt.Sprint(expected)
		if a != b {
			t.Errorf("#%d: failed, got %v, expected %v", i, a, b)
		}
	}

	check(-1)
	for i := 0; i < 2; i++ {
		<-c
		check(i, i)
	}
}
