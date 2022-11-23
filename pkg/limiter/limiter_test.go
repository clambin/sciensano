package limiter_test

import (
	"bytes"
	"context"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/pkg/limiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	const maxParallel = 3
	c := &Caller{}
	l := limiter.NewLimiter(c, maxParallel)

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	_, err := l.Do(req)
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req2, _ := http.NewRequest(http.MethodGet, "", nil)
			_, err2 := l.Do(req2)
			require.NoError(t, err2)
		}()
	}
	wg.Wait()
	assert.LessOrEqual(t, c.GetMax(), maxParallel)
}

/*
TODO: this sometimes fails: no error is received?
func TestLimiter_Timeout(t *testing.T) {
	const maxParallel = 3
	c := &Caller{Delay: time.Minute}
	l := limiter.NewLimiter(c, maxParallel)

	for i := 0; i < 50; i++ {
		go func() {
			req, _ := http.NewRequest(http.MethodGet, "", nil)
			_, _ = l.Do(req)
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	_, err := l.Do(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
	cancel()
}
*/

type Caller struct {
	Delay   time.Duration
	lock    sync.RWMutex
	current int
	max     int
}

var _ httpclient.Caller = &Caller{}

func (c *Caller) Do(req *http.Request) (*http.Response, error) {
	c.wait(req.Context())
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}, nil
}

func (c *Caller) getDelay() time.Duration {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Delay == 0 {
		c.Delay = 10 * time.Millisecond
	}
	return c.Delay
}

func (c *Caller) increase() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.current++
	if c.current > c.max {
		c.max = c.current
	}
}

func (c *Caller) decrease() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.current--
}

func (c *Caller) GetMax() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.max
}

func (c *Caller) wait(ctx context.Context) {
	c.increase()
	select {
	case <-ctx.Done():
		break
	case <-time.After(c.getDelay()):
		break
	}
	c.decrease()
}
