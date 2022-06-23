package fetcher_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	const maxParallel = 3
	c := &Caller{}
	l := fetcher.NewLimiter(c, maxParallel)

	ctx := context.Background()

	_, err := l.Fetch(ctx, 0)
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			_, err2 := l.GetLastUpdates(ctx, 0)
			require.NoError(t, err2)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			_, err2 := l.Fetch(ctx, 0)
			require.NoError(t, err2)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.LessOrEqual(t, c.GetMax(), maxParallel)
}

func TestLimiter_Timeout(t *testing.T) {
	c := &Caller{}
	l := fetcher.NewLimiter(c, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	_, err := l.Fetch(ctx, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
	_, err = l.GetLastUpdates(ctx, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
	cancel()
}

type Caller struct {
	Delay   time.Duration
	lock    sync.RWMutex
	current int
	max     int
}

var _ fetcher.Fetcher = &Caller{}

func (c *Caller) Fetch(ctx context.Context, _ int) (results []apiclient.APIResponse, err error) {
	c.wait(ctx)
	return
}

func (c *Caller) GetLastUpdates(ctx context.Context, _ int) (lastModified time.Time, err error) {
	c.wait(ctx)
	return
}

func (c *Caller) DataTypes() map[int]string {
	return map[int]string{
		0: "test",
	}
}

func (c *Caller) getDelay() time.Duration {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Delay == 0 {
		c.Delay = 50 * time.Millisecond
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
