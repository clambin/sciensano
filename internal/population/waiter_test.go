package population_test

import (
	"context"
	"github.com/clambin/sciensano/v2/internal/population"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWaiter_WaitTillReady(t *testing.T) {
	testcases := []struct {
		name    string
		delay   time.Duration
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "ready",
			delay:   100 * time.Millisecond,
			wantErr: assert.NoError,
		},
		{
			name:    "not ready",
			delay:   time.Second,
			wantErr: assert.Error,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			w := population.Waiter{}
			go func() {
				time.Sleep(tt.delay)
				w.Ready()
			}()

			tt.wantErr(t, w.WaitTillReady(ctx))
		})
	}
}

func TestWaiter_AlreadyReady(t *testing.T) {
	w := population.Waiter{}
	w.Ready()

	assert.NoError(t, w.WaitTillReady(context.Background()))
}
