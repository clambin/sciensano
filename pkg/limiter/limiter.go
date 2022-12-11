package limiter

import (
	"golang.org/x/sync/semaphore"
	"net/http"
)

type Limiter struct {
	r     http.RoundTripper
	limit *semaphore.Weighted
}

var _ http.RoundTripper = &Limiter{}

func NewLimiter(caller http.RoundTripper, maxParallel int64) *Limiter {
	return &Limiter{
		r:     caller,
		limit: semaphore.NewWeighted(maxParallel),
	}
}

func (l *Limiter) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	err = l.limit.Acquire(req.Context(), 1)
	if err != nil {
		return
	}

	defer l.limit.Release(1)
	return l.r.RoundTrip(req)
}
