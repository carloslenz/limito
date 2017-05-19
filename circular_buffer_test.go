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
			if _, err := cb.add(); err != nil {
				t.Errorf("#%d failed, got %q, expected nil", i, err)
			}
		}
		if _, err := cb.add(); err != errCircularBufferFull {
			t.Errorf("#%d failed, got %q expected %q", i, err, errCircularBufferFull)
		}
	})
	t.Run("Complex", func(t *testing.T) {
		cb := circularBuffer{buf: make([]chan struct{}, 4)}

		var (
			sizes   []int
			count   int
			errors  []string
			indexes [][]int
		)
		ids := make(map[chan struct{}][]int, 4)
		consume := func(c chan struct{}) {
			i := ids[c][0]
			ids[c] = ids[c][1:]
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
		inputs := 9
		for inputs > 0 {
			putIndexes()
			c, err := cb.add()
			if gotError(err) {
				removeAll()
			} else {
				ids[c] = append(ids[c], 9-inputs)
				inputs--
				count++
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
			{0, 0},                                 // removing (failed)
			{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, // adding (last failed)
			{0, 4}, {1, 4}, {2, 4}, {3, 4}, {0, 0}, // removing (last failed)
			{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, // adding (last failed)
			{0, 4}, {1, 4}, {2, 4}, {3, 4}, {0, 0}, // removing (last failed)
			{0, 0}, {0, 1}, // adding
			{0, 0}, // removing
		}
		if !reflect.DeepEqual(indexes, expectedIndexes) {
			t.Errorf("indexes: failed, got %#v, expected %#v", indexes, expectedIndexes)
		}
	})
}
