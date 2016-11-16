package limito

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestWaitList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wl := NewWaitList(2)
	results := make(chan string, 3)

	const N = 3
	var wg sync.WaitGroup
	wg.Add(N - 1)

	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			if err := wl.Wait(ctx); err != nil {
				results <- fmt.Sprint(err)
			} else {
				results <- ""
			}
		}(i)
	}

	wg.Wait()

	wg.Add(1)
	cancel()
	wg.Wait()

	close(results)
	var obtained []string
	for s := range results {
		obtained = append(obtained, s)
	}
	expected := []string{
		"",
		"",
		fmt.Sprint(ctx.Err()),
	}
	if !reflect.DeepEqual(obtained, expected) {
		t.Errorf("failed, got %s, expected %s", obtained, expected)
	}
}
