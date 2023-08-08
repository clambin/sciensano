package population

import (
	"context"
	"sync"
	"time"
)

type Waiter struct {
	Interval time.Duration
	ready    bool
	lock     sync.RWMutex
}

func (w *Waiter) isReady() bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.ready
}

func (w *Waiter) Ready() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.ready = true
}

func (w *Waiter) WaitTillReady(ctx context.Context) error {
	interval := w.Interval
	if interval == 0 {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if w.isReady() {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
