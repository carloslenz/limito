package limito

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

type middleware struct {
	http.Handler
	sync.Mutex
	n        int
	limiters map[string]*Limiter
}

// Middleware wraps h with Limiter of up to n concurrrent requests per ID, after which
// http.StatusTooManyRequests (429) is returned to the client. If the ID is not set in
// the request's context, it panics.
func Middleware(n int, h http.Handler) http.Handler {
	return &middleware{
		Handler:  h,
		n:        n,
		limiters: make(map[string]*Limiter),
	}
}

const (
	tooManyConcurrentRequests = "too many concurrent requests"
)

func (m *middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	id := GetMiddlewareID(r.Context())
	lim, err := m.limit(id)
	if err != nil {
		http.Error(rw, tooManyConcurrentRequests, http.StatusTooManyRequests)
		return
	}
	defer lim.Leave()
	m.Handler.ServeHTTP(rw, r)
}

func (m *middleware) limit(id string) (*Limiter, error) {
	m.Mutex.Lock()

	lim := m.limiters[id]
	if lim == nil {
		l := NewLimiter(m.n)
		lim = &l
		m.limiters[id] = lim
	}
	err := lim.Enter()

	m.Mutex.Unlock()
	return lim, err
}

type middlewareKey int

var (
	limitoMiddlewareID middlewareKey = 1

	errMiddlewareIDNotFound = errors.New("limito middleware ID: not found")
)

// GetMiddlewareID extracts ID from ctx. If there is ID, it panics.
func GetMiddlewareID(ctx context.Context) string {
	if id, ok := ctx.Value(limitoMiddlewareID).(string); ok {
		return id
	}
	// Yes panic, because it is a bug if the ID is missing, not an error!
	panic(errMiddlewareIDNotFound)
}

// SetMiddlewareID stores ID in a new context.
func SetMiddlewareID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, limitoMiddlewareID, id)
}
