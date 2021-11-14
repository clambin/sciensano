package measurement_test

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/measurement"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type CacheTester struct {
	measurement.Cache
	Count int
}

func (ct *CacheTester) set(_ context.Context) (value []measurement.Measurement, err error) {
	ct.Count++
	return
}

func (ct *CacheTester) fail(_ context.Context) (value []measurement.Measurement, err error) {
	return nil, fmt.Errorf("failed")
}

func TestCache(t *testing.T) {
	ctx := context.Background()
	ct := &CacheTester{
		Cache: measurement.Cache{
			Retention: 15 * time.Minute,
		},
	}
	assert.Zero(t, ct.CacheSize())

	_, err := ct.Call(ctx, "test", ct.fail)
	assert.Error(t, err)
	assert.Equal(t, 1, ct.CacheSize())

	_, err = ct.Call(ctx, "test", ct.set)
	assert.NoError(t, err)
	assert.Equal(t, 1, ct.CacheSize())
	assert.Equal(t, 1, ct.Count)

	_, err = ct.Call(ctx, "test", ct.set)
	assert.Equal(t, 1, ct.CacheSize())
	assert.Equal(t, 1, ct.Count)

	_, err = ct.Call(ctx, "test", ct.fail)
	assert.NoError(t, err)
	assert.Equal(t, 1, ct.CacheSize())
	assert.Equal(t, 1, ct.Count)
}
