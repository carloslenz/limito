package limito

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCircularBuffer(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		cb := circularBuffer{buf: make([]chan struct{}, 2)}
		i := 0
		for ; i < 2; i++ {
			if err := cb.add(nil); err != nil {
				t.Errorf("#%d failed, got %q, expected nil", i, err)
			}
		}
		if err := cb.add(nil); err != errCircularBufferFull {
			t.Errorf("#%d failed, got %q expected %q", i, err, errCircularBufferFull)
		}
	})
	t.Run("Complex", func(t *testing.T) {
		inputs := make([]chan struct{}, 9)
		for i := 0; i < len(inputs); i++ {
			c := make(chan struct{}, i)
			for j := 0; j < i; j++ {
				c <- struct{}{}
			}
			close(c)
			inputs[i] = c
		}
		cb := circularBuffer{buf: make([]chan struct{}, 4)}

		var (
			sizes   []int
			count   int
			errors  []string
			indexes [][]int
		)
		consume := func(c chan struct{}) {
			var i int
			for _ = range c {
				i++
			}
			sizes = append(sizes, i)
		}
		fmtError := func(err error, count int) string {
			return fmt.Sprintf("%s: %d", err, count)
		}
		gotError := func(err error) bool {
			if err != nil {
				errors = append(errors, fmtError(err, count))
				return true
			}
			return false
		}
		putIndexes := func() {
			indexes = append(indexes, []int{cb.read, cb.write})
		}
		removeAll := func() {
			for {
				putIndexes()
				c, err := cb.remove()
				if gotError(err) {
					break
				} else {
					consume(c)
				}
			}
		}

		removeAll()
		for len(inputs) > 0 {
			putIndexes()
			if gotError(cb.add(inputs[0])) {
				removeAll()
			} else {
				inputs, count = inputs[1:], count+1
			}
		}
		removeAll()

		expectedErrors := []string{
			fmtError(errCircularBufferEmpty, 0),
			fmtError(errCircularBufferFull, 4),
			fmtError(errCircularBufferEmpty, 4),
			fmtError(errCircularBufferFull, 8),
			fmtError(errCircularBufferEmpty, 8),
			fmtError(errCircularBufferEmpty, 9),
		}
		if !reflect.DeepEqual(errors, expectedErrors) {
			t.Errorf("errors: failed, got %#v, expected %#v", errors, expectedErrors)
		}

		expectedSizes := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
		if !reflect.DeepEqual(sizes, expectedSizes) {
			t.Errorf("sizes: failed, got %#v, expected %#v", sizes, expectedSizes)
		}

		expectedIndexes := [][]int{
			[]int{0, 0},                                                     // removing (failed)
			[]int{0, 0}, []int{0, 1}, []int{0, 2}, []int{0, 3}, []int{0, 4}, // adding (last failed)
			[]int{0, 4}, []int{1, 4}, []int{2, 4}, []int{3, 4}, []int{0, 0}, // removing (last failed)
			[]int{0, 0}, []int{0, 1}, []int{0, 2}, []int{0, 3}, []int{0, 4}, // adding (last failed)
			[]int{0, 4}, []int{1, 4}, []int{2, 4}, []int{3, 4}, []int{0, 0}, // removing (last failed)
			[]int{0, 0}, []int{0, 1}, // adding
			[]int{0, 0}, // removing
		}
		if !reflect.DeepEqual(indexes, expectedIndexes) {
			t.Errorf("indexes: failed, got %#v, expected %#v", indexes, expectedIndexes)
		}
	})
}
