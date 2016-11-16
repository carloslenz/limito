package limito

import (
	"fmt"
	"testing"
)

func TestLimiter(t *testing.T) {
	ok := toString(nil)
	empty := toString(errLimiterEmpty)
	full := toString(errLimiterFull)
	lim := NewLimiter(2)
	defer lim.Close()
	tests := []struct {
		got      error
		expected string
	}{
		{lim.Leave(), empty},
		{lim.Enter(), ok},
		{lim.Enter(), ok},
		{lim.Enter(), full},
		{lim.Enter(), full},
		{lim.Leave(), ok},
		{lim.Leave(), ok},
		{lim.Enter(), ok},
		{lim.Leave(), ok},
		{lim.Leave(), empty},
	}
	for i, test := range tests {
		obtained := toString(test.got)
		if obtained != test.expected {
			t.Errorf("%d: failed, got %q, expected %q", i, obtained, test.expected)
		}
	}
	lim.Close()
	// Should not panic.
	lim.Close()
}

func toString(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprint(err)
}
