package limiter

import (
	"github.com/clambin/httpclient"
	"golang.org/x/sync/semaphore"
	"net/http"
)

type Limiter struct {
	httpclient.Caller
	limit *semaphore.Weighted
}

var _ httpclient.Caller = &Limiter{}

func NewLimiter(caller httpclient.Caller, maxParallel int64) *Limiter {
	return &Limiter{
		Caller: caller,
		limit:  semaphore.NewWeighted(maxParallel),
	}
}

func (l *Limiter) Do(req *http.Request) (resp *http.Response, err error) {
	err = l.limit.Acquire(req.Context(), 1)
	if err != nil {
		return
	}

	defer l.limit.Release(1)
	return l.Caller.Do(req)
}
